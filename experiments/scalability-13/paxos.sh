arrival=$1
iteration=$2

conflicts=2
replicaBatchSize=2000
replicaBatchTime=3000

pwd=$(pwd)
. "${pwd}"/experiments/setup-13/ip.sh

remote_algo_path="/raxos/binary/epaxos_server"
remote_ctl_path="/raxos/binary/epaxos_client"
remote_master_path="/raxos/binary/epaxos_master"
remote_log_path="/home/${user_name}/raxos/logs/"

echo "Starting test"

output_path="${pwd}/experiments/scalability-13/logs/paxos/${arrival}/${iteration}/"
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

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1} ".${remote_master_path} -N 13 " >${output_path}1.log &

sleep 5

echo "starting replicas"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica1}  ".${remote_algo_path}  -port 10070 -maddr ${replica1_ip} -addr ${replica1_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}2.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica2}  ".${remote_algo_path}  -port 10071 -maddr ${replica1_ip} -addr ${replica2_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}3.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica3}  ".${remote_algo_path}  -port 10072 -maddr ${replica1_ip} -addr ${replica3_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}4.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica4}  ".${remote_algo_path}  -port 10073 -maddr ${replica1_ip} -addr ${replica4_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}5.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica5}  ".${remote_algo_path}  -port 10074 -maddr ${replica1_ip} -addr ${replica5_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}6.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica6}  ".${remote_algo_path}  -port 10075 -maddr ${replica1_ip} -addr ${replica6_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}7.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica7}  ".${remote_algo_path}  -port 10076 -maddr ${replica1_ip} -addr ${replica7_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}8.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica8}  ".${remote_algo_path}  -port 10077 -maddr ${replica1_ip} -addr ${replica8_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}9.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica9}  ".${remote_algo_path}  -port 10078 -maddr ${replica1_ip} -addr ${replica9_ip}  -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}10.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica10} ".${remote_algo_path}  -port 10079 -maddr ${replica1_ip} -addr ${replica10_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}11.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica11} ".${remote_algo_path}  -port 10080 -maddr ${replica1_ip} -addr ${replica11_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}12.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica12} ".${remote_algo_path}  -port 10081 -maddr ${replica1_ip} -addr ${replica12_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}13.log &
sleep 5
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} -n -f ${replica13} ".${remote_algo_path}  -port 10082 -maddr ${replica1_ip} -addr ${replica13_ip} -batchSize ${replicaBatchSize} -batchTime ${replicaBatchTime}     -exec  -dreply -pa" >${output_path}14.log &
sleep 5

echo "Starting client[s]"

nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client1}  ".${remote_ctl_path} -name 15    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 0 -logFilePath ${remote_log_path} -l" >${output_path}15.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client2}  ".${remote_ctl_path} -name 16    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 1 -logFilePath ${remote_log_path} -l" >${output_path}16.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client3}  ".${remote_ctl_path} -name 17    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 2 -logFilePath ${remote_log_path} -l" >${output_path}17.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client4}  ".${remote_ctl_path} -name 18    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 3 -logFilePath ${remote_log_path} -l" >${output_path}18.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client5}  ".${remote_ctl_path} -name 19    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 4 -logFilePath ${remote_log_path} -l" >${output_path}19.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client6}  ".${remote_ctl_path} -name 20    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 5 -logFilePath ${remote_log_path} -l" >${output_path}20.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client7}  ".${remote_ctl_path} -name 21    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 6 -logFilePath ${remote_log_path} -l" >${output_path}21.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client8}  ".${remote_ctl_path} -name 22    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 2 -logFilePath ${remote_log_path} -l" >${output_path}22.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client9}  ".${remote_ctl_path} -name 23    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 3 -logFilePath ${remote_log_path} -l" >${output_path}23.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client10} ".${remote_ctl_path} -name 24    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 4 -logFilePath ${remote_log_path} -l" >${output_path}24.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client11} ".${remote_ctl_path} -name 25    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 0 -logFilePath ${remote_log_path} -l" >${output_path}25.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client12} ".${remote_ctl_path} -name 26    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 1 -logFilePath ${remote_log_path} -l" >${output_path}26.log &
nohup sshpass ssh -o "StrictHostKeyChecking no" -i ${cert} ${client13} ".${remote_ctl_path} -name 27    -maddr ${replica1_ip} -w 50  -c ${conflicts} -arrivalRate ${arrival} -clientBatchSize 50 -defaultReplica 2 -logFilePath ${remote_log_path} -l" >${output_path}27.log &


sleep 100

echo "Completed Client[s]"

sleep 10

echo "finished test"

echo "Finish test"