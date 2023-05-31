arrival=$1
replicaBatchSize=$2
replicaBatchTime=$3
clientBatchSize=$4
pipelineLength=$5
iteration=$6

pwd=$(pwd)
. "${pwd}"/aws/best-case-lan/setup/ip.sh

remote_algo_path="/raxos/bin/epaxos_server"
remote_ctl_path="/raxos/bin/epaxos_client"
remote_master_path="/raxos/bin/epaxos_master"
remote_log_path="/home/ubuntu/raxos/logs/"
output_path="${pwd}/aws/best-case-lan/epaxos_paxos/logs/paxos/${iteration}/${arrival}/${replicaBatchSize}/${replicaBatchTime}/${clientBatchSize}/${pipelineLength}/"
rm -r "${output_path}" ; mkdir -p "${output_path}"

pwd=$(pwd)
. "${pwd}"/aws/best-case-lan/setup/ip.sh

reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"
kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill paxos_raft_repl ; pkill paxos_raft_clie; pkill raxos_replica ; pkill raxos_client"

sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica1} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica2} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica3} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica4} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica5} "${reset_logs}; ${kill_command}"

sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} "${reset_logs}; ${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} "${reset_logs}; ${kill_command}"

sleep 5
echo "Removed all files in remote servers"

sleep 2

echo "starting master"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_master_path} -N 5 " >${output_path}1.log &

sleep 5

echo "starting replicas"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path}  -port 10070 -maddr ${replica1_ip} -addr ${replica1_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength} -exec  -dreply -pa" >${output_path}2.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path}  -port 10071 -maddr ${replica1_ip} -addr ${replica2_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength} -exec  -dreply -pa" >${output_path}3.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path}  -port 10072 -maddr ${replica1_ip} -addr ${replica3_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength} -exec  -dreply -pa" >${output_path}4.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4} ".${remote_algo_path}  -port 10073 -maddr ${replica1_ip} -addr ${replica4_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength} -exec  -dreply -pa" >${output_path}5.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5} ".${remote_algo_path}  -port 10074 -maddr ${replica1_ip} -addr ${replica5_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength} -exec  -dreply -pa" >${output_path}6.log &
sleep 5

sleep 5

echo "Starting client[s]"

nohup  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} -name 7    -maddr  ${replica1_ip} -w 50  -c 2 -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 0 -l -logFilePath ${remote_log_path} " >${output_path}7.log &
nohup  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} -name 8    -maddr  ${replica1_ip} -w 50  -c 2 -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 1 -l -logFilePath ${remote_log_path} " >${output_path}8.log &
nohup  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path} -name 9    -maddr  ${replica1_ip} -w 50  -c 2 -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 2 -l -logFilePath ${remote_log_path} " >${output_path}9.log &
nohup  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} ".${remote_ctl_path} -name 10   -maddr  ${replica1_ip} -w 50  -c 2 -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 3 -l -logFilePath ${remote_log_path} " >${output_path}10.log &
nohup  sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} ".${remote_ctl_path} -name 11   -maddr  ${replica1_ip} -w 50  -c 2 -arrivalRate ${arrival} -clientBatchSize ${clientBatchSize} -defaultReplica 4 -l -logFilePath ${remote_log_path} " >${output_path}11.log &

sleep 100

echo "Completed Client[s]"

sleep 10

sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica1} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica2} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica3} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica4} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${replica5} "${kill_command}"

sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} "${kill_command}"
sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} "${kill_command}"

echo "killed  instances"

echo "Finish test"