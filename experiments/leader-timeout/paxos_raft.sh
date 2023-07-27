algo=$1
iteration=$2
timeout=$3
mode=$4 # performance or recovery

arrival=5000
replicaBatchSize=3000
replicaBatchTime=5000

pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

remote_algo_path="/raxos/binary/paxos_raft_repl"
remote_ctl_path="/raxos/binary/paxos_raft_clie"
remote_log_path="/home/${user_name}/raxos/logs/"

remote_config_path="/home/${user_name}/raxos/binary/paxos.yml"

echo "Starting test"

output_path="${pwd}/experiments/leader-timeout/logs/${mode}/${algo}/${iteration}/${timeout}/"
rm -r "${output_path}" ; mkdir -p "${output_path}"

echo "Removed old local log files"

reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"
kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill paxos_raft_repl ; pkill paxos_raft_clie; pkill quepaxa_replica ; pkill quepaxa_client; pkill rabia"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}; ${kill_command}"
done

sleep 5
echo "Removed all files in remote servers"

sleep 2

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path} --name 1 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --consAlgo ${algo}  --viewTimeout ${timeout} " >${output_path}1.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path} --name 2 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --consAlgo ${algo}  --viewTimeout ${timeout} " >${output_path}2.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path} --name 3 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --consAlgo ${algo}  --viewTimeout ${timeout} " >${output_path}3.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4} ".${remote_algo_path} --name 4 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --consAlgo ${algo}  --viewTimeout ${timeout} " >${output_path}4.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5} ".${remote_algo_path} --name 5 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --consAlgo ${algo}  --viewTimeout ${timeout} " >${output_path}5.log &
echo "Started servers"

sleep 10

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 1" >${output_path}status1.log &
echo "Sent initial status"

sleep 5

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 3" >${output_path}status3.log &
echo "Sent consensus start up"

sleep 15

echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500 --window 10000" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500 --window 10000" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path} --name 23 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500 --window 10000" >${output_path}23.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} ".${remote_ctl_path} --name 24 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500 --window 10000" >${output_path}24.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} ".${remote_ctl_path} --name 25 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500 --window 10000" >${output_path}25.log &

sleep 15

if [[ "${mode}" == "recovery" ]]
then
  echo "killing the first leader"
  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "${replica1}" "${kill_command}"
fi

sleep 110

echo "Completed Client[s]"

if [[ "${mode}" == "recovery" ]]
then
  scp -i ${cert} ${client1}:${remote_log_path}21.txt ${output_path}21.txt
  scp -i ${cert} ${client2}:${remote_log_path}22.txt ${output_path}22.txt
  scp -i ${cert} ${client3}:${remote_log_path}23.txt ${output_path}23.txt
  scp -i ${cert} ${client4}:${remote_log_path}24.txt ${output_path}24.txt
  scp -i ${cert} ${client5}:${remote_log_path}25.txt ${output_path}25.txt

fi


echo "Finish test"
