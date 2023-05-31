arrival=$1
replicaBatchSize=$2
replicaBatchTime=$3
setting=$4 # LAN or WAN
iteration=$5

conflicts=2
pipelineLength=1

if [[ "${setting}" == "WAN" ]]
then
  pipelineLength=20
fi

pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

remote_algo_path="/raxos/binary/epaxos_server"
remote_ctl_path="/raxos/binary/epaxos_client"
remote_master_path="/raxos/binary/epaxos_master"
remote_log_path="/home/${user_name}/raxos/logs/"

echo "Starting execution latency test"

output_path="${pwd}/experiments/best-case/logs/epaxos/${arrival}/${replicaBatchSize}/${replicaBatchTime}/${setting}/${iteration}/execution/"
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

echo "starting master"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_master_path} -N 5 " >${output_path}1.log &

sleep 5

echo "starting replicas"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path}  -port 10070 -maddr ${replica1_ip} -addr ${replica1_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -exec  -dreply -e" >${output_path}2.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path}  -port 10071 -maddr ${replica1_ip} -addr ${replica2_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -exec  -dreply -e" >${output_path}3.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path}  -port 10072 -maddr ${replica1_ip} -addr ${replica3_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -exec  -dreply -e" >${output_path}4.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4} ".${remote_algo_path}  -port 10073 -maddr ${replica1_ip} -addr ${replica4_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -exec  -dreply -e" >${output_path}5.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5} ".${remote_algo_path}  -port 10074 -maddr ${replica1_ip} -addr ${replica5_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -exec  -dreply -e" >${output_path}6.log &
sleep 5

echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} -name 7    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 0 -logFilePath ${remote_log_path} " >${output_path}7.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} -name 8    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 1 -logFilePath ${remote_log_path} " >${output_path}8.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path} -name 9    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 2 -logFilePath ${remote_log_path} " >${output_path}9.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} ".${remote_ctl_path} -name 10   -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 3 -logFilePath ${remote_log_path} " >${output_path}10.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} ".${remote_ctl_path} -name 11   -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 4 -logFilePath ${remote_log_path} " >${output_path}11.log &

sleep 100

echo "Completed Client[s]"

sleep 10

echo "finished execution latency test"
echo "starting commit latency test"

output_path="${pwd}/experiments/best-case/logs/epaxos/${arrival}/${replicaBatchSize}/${replicaBatchTime}/${setting}/${iteration}/commit/"
rm -r "${output_path}" ; mkdir -p "${output_path}"

echo "Removed old local log files"

reset_logs="rm -r ${remote_log_path} ; mkdir -p ${remote_log_path}"
kill_command="pkill epaxos_master ; pkill epaxos_server; pkill epaxos_client; pkill paxos_raft_repl ; pkill paxos_raft_clie; pkill quepaxa_replica ; pkill quepaxa_client"

for i in "${machines[@]}"
do
   echo "killing instances and removing old files ${i}"
   sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} "$i" "${reset_logs}; ${kill_command}"
done

sleep 5
echo "Removed all files in remote servers"

sleep 2

echo "starting master"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_master_path} -N 5 " >${output_path}1.log &

sleep 5

echo "starting replicas"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_algo_path}  -port 10070 -maddr ${replica1_ip} -addr ${replica1_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -e" >${output_path}2.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2} ".${remote_algo_path}  -port 10071 -maddr ${replica1_ip} -addr ${replica2_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -e" >${output_path}3.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3} ".${remote_algo_path}  -port 10072 -maddr ${replica1_ip} -addr ${replica3_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -e" >${output_path}4.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4} ".${remote_algo_path}  -port 10073 -maddr ${replica1_ip} -addr ${replica4_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -e" >${output_path}5.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5} ".${remote_algo_path}  -port 10074 -maddr ${replica1_ip} -addr ${replica5_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}  -pipeline ${pipelineLength}   -e" >${output_path}6.log &
sleep 5

echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1} ".${remote_ctl_path} -name 7    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 0 -logFilePath ${remote_log_path} " >${output_path}7.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2} ".${remote_ctl_path} -name 8    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 1 -logFilePath ${remote_log_path} " >${output_path}8.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3} ".${remote_ctl_path} -name 9    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 2 -logFilePath ${remote_log_path} " >${output_path}9.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4} ".${remote_ctl_path} -name 10   -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 3 -logFilePath ${remote_log_path} " >${output_path}10.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5} ".${remote_ctl_path} -name 11   -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 4 -logFilePath ${remote_log_path} " >${output_path}11.log &

sleep 100

echo "Completed Client[s]"

echo "Finish test"