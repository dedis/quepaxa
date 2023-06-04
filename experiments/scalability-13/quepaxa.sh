arrival=$1
iteration=$2

replicaBatchSize=2000
replicaBatchTime=3000

pwd=$(pwd)
. "${pwd}"/experiments/setup-13/ip.sh

remote_algo_path="/raxos/binary/quepaxa_replica"
remote_ctl_path="/raxos/binary/quepaxa_client"
remote_log_path="/home/${user_name}/raxos/logs/"

remote_config_path="/home/${user_name}/raxos/binary/quepaxa.yml"

echo "Starting test"

output_path="${pwd}/experiments/scalability-13/logs/quepaxa/${arrival}/${iteration}/"
rm -r "${output_path}" ; mkdir -p "${output_path}"

echo "Removed old local log files"

reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"
kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill paxos_raft_repl ; pkill paxos_raft_clie; pkill quepaxa_replica ; pkill quepaxa_client"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files in ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}"
done

sleep 5
echo "Removed all files in remote servers"

sleep 2

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1}  ".${remote_algo_path} --name 1  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}1.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2}  ".${remote_algo_path} --name 2  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}2.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3}  ".${remote_algo_path} --name 3  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}3.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4}  ".${remote_algo_path} --name 4  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}4.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5}  ".${remote_algo_path} --name 5  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}5.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica6}  ".${remote_algo_path} --name 6  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}6.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica7}  ".${remote_algo_path} --name 7  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}7.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica8}  ".${remote_algo_path} --name 8  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}8.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica9}  ".${remote_algo_path} --name 9  --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}9.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica10} ".${remote_algo_path} --name 10 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}10.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica11} ".${remote_algo_path} --name 11 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}11.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica12} ".${remote_algo_path} --name 12 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}12.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica13} ".${remote_algo_path} --name 13 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} " >${output_path}13.log &
echo "Started servers"

sleep 5

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status  --operationType 1" >"${output_path}"status1.log &
echo "Sent initial status"

sleep 35

echo "Starting client"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1}  ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2}  ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3}  ".${remote_ctl_path} --name 23 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}23.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4}  ".${remote_ctl_path} --name 24 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}24.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5}  ".${remote_ctl_path} --name 25 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}25.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client6}  ".${remote_ctl_path} --name 26 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}26.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client7}  ".${remote_ctl_path} --name 27 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}27.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client8}  ".${remote_ctl_path} --name 28 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}28.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client9}  ".${remote_ctl_path} --name 29 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}29.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client10} ".${remote_ctl_path} --name 30 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}30.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client11} ".${remote_ctl_path} --name 31 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}31.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client12} ".${remote_ctl_path} --name 32 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}32.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client13} ".${remote_ctl_path} --name 33 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}33.log &

sleep 110

echo "Completed Client[s]"

echo "finished test"

echo "Finish test"