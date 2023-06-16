arrival=$1
replicaBatchSize=$2
replicaBatchTime=$3
setting=$4 # LAN or WAN
iteration=$5
serverMode=$6
requestPropagationTime=$7

pipelineLength=1

if [[ "${setting}" == "WAN" ]]
then
  pipelineLength=10
fi

pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

remote_algo_path="/raxos/binary/quepaxa_replica"
remote_ctl_path="/raxos/binary/quepaxa_client"
remote_log_path="/home/${user_name}/raxos/logs/"

remote_config_path="/home/${user_name}/raxos/binary/quepaxa.yml"

echo "Starting execution latency test"

output_path="${pwd}/experiments/best-case/logs/quepaxa/${arrival}/${replicaBatchSize}/${replicaBatchTime}/${setting}/${iteration}/${serverMode}/${requestPropagationTime}/execution/"
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

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path} --name 1 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --pipelineLength ${pipelineLength} --serverMode ${serverMode} --requestPropagationTime ${requestPropagationTime} " >${output_path}1.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path} --name 2 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --pipelineLength ${pipelineLength} --serverMode ${serverMode} --requestPropagationTime ${requestPropagationTime} " >${output_path}2.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path} --name 3 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --pipelineLength ${pipelineLength} --serverMode ${serverMode} --requestPropagationTime ${requestPropagationTime} " >${output_path}3.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4} ".${remote_algo_path} --name 4 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --pipelineLength ${pipelineLength} --serverMode ${serverMode} --requestPropagationTime ${requestPropagationTime} " >${output_path}4.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5} ".${remote_algo_path} --name 5 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize ${replicaBatchSize} --batchTime ${replicaBatchTime} --pipelineLength ${pipelineLength} --serverMode ${serverMode} --requestPropagationTime ${requestPropagationTime} " >${output_path}5.log &
echo "Started servers"

sleep 5

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status  --operationType 1" >"${output_path}"status1.log &
echo "Sent initial status"

sleep 35

echo "Starting client"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path} --name 23 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}23.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} ".${remote_ctl_path} --name 24 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}24.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} ".${remote_ctl_path} --name 25 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrival}  --batchSize 50 --batchTime 500" >${output_path}25.log &

sleep 110

echo "Completed Client[s]"

echo "finished execution latency test"

echo "Finish test"