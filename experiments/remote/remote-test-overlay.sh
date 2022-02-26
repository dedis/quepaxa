arrivalRate=$1

remote_algo_path="/raxos/temp/replica/bin/replica"
remote_ctl_path="/raxos/temp/client/bin/client"
remote_log_path="/home/pasindu/raxos/temp/logs/"
remote_config_path="/home/pasindu/raxos/temp/configuration/remote/configuration.yml"
output_path="logs/"

replica1=pasindu@dedis-140.icsil1.epfl.ch
replica1_cert="/home/pasindu/Pictures/pasindu_rsa"
replica2=pasindu@dedis-141.icsil1.epfl.ch
replica2_cert="/home/pasindu/Pictures/pasindu_rsa"
replica3=pasindu@dedis-142.icsil1.epfl.ch
replica3_cert="/home/pasindu/Pictures/pasindu_rsa"
replica4=pasindu@dedis-143.icsil1.epfl.ch
replica4_cert="/home/pasindu/Pictures/pasindu_rsa"
replica5=pasindu@dedis-144.icsil1.epfl.ch
replica5_cert="/home/pasindu/Pictures/pasindu_rsa"

client1=pasindu@dedis-145.icsil1.epfl.ch
client1_cert="/home/pasindu/Pictures/pasindu_rsa"
client2=pasindu@dedis-146.icsil1.epfl.ch
client2_cert="/home/pasindu/Pictures/pasindu_rsa"
client3=pasindu@dedis-147.icsil1.epfl.ch
client3_cert="/home/pasindu/Pictures/pasindu_rsa"
client4=pasindu@dedis-148.icsil1.epfl.ch
client4_cert="/home/pasindu/Pictures/pasindu_rsa"
client5=pasindu@dedis-149.icsil1.epfl.ch
client5_cert="/home/pasindu/Pictures/pasindu_rsa"

rm ${output_path}0.txt
rm ${output_path}0.log
rm ${output_path}1.txt
rm ${output_path}1.log
rm ${output_path}2.txt
rm ${output_path}2.log
rm ${output_path}3.txt
rm ${output_path}3.log
rm ${output_path}4.txt
rm ${output_path}4.log
rm ${output_path}5.txt
rm ${output_path}5.log
rm ${output_path}6.txt
rm ${output_path}6.log
rm ${output_path}7.txt
rm ${output_path}7.log
rm ${output_path}8.txt
rm ${output_path}8.log
rm ${output_path}9.txt
rm ${output_path}9.log

rm ${output_path}status1.log
rm ${output_path}status2.log
rm ${output_path}pythonreport.log
rm ${output_path}test.log

echo "Removed old log files"

sshpass ssh -i ${replica1_cert} ${replica1} "rm ${remote_log_path}0.txt"
sshpass ssh -i ${replica2_cert} ${replica2} "rm ${remote_log_path}1.txt"
sshpass ssh -i ${replica3_cert} ${replica3} "rm ${remote_log_path}2.txt"
sshpass ssh -i ${replica4_cert} ${replica4} "rm ${remote_log_path}3.txt"
sshpass ssh -i ${replica5_cert} ${replica5} "rm ${remote_log_path}4.txt"

sshpass ssh -i ${client1_cert} ${client1} "rm ${remote_log_path}5.txt"
sshpass ssh -i ${client2_cert} ${client2} "rm ${remote_log_path}6.txt"
sshpass ssh -i ${client3_cert} ${client3} "rm ${remote_log_path}7.txt"
sshpass ssh -i ${client4_cert} ${client4} "rm ${remote_log_path}8.txt"
sshpass ssh -i ${client5_cert} ${client5} "rm ${remote_log_path}9.txt"

sleep 5
echo "Removed all files in remote servers"

kill_command="pkill replica ; pkill client"

sshpass ssh -i ${replica1_cert} ${replica1} "${kill_command}"
sshpass ssh -i ${replica2_cert} ${replica2} "${kill_command}"
sshpass ssh -i ${replica3_cert} ${replica3} "${kill_command}"
sshpass ssh -i ${replica4_cert} ${replica4} "${kill_command}"
sshpass ssh -i ${replica5_cert} ${replica5} "${kill_command}"

sshpass ssh -i ${client1_cert} ${client1} "${kill_command}"
sshpass ssh -i ${client2_cert} ${client2} "${kill_command}"
sshpass ssh -i ${client3_cert} ${client3} "${kill_command}"
sshpass ssh -i ${client4_cert} ${client4} "${kill_command}"
sshpass ssh -i ${client5_cert} ${client5} "${kill_command}"

echo "killed previous running instances"

sleep 5

nohup sshpass ssh -i ${replica1_cert} -n -f ${replica1} ".${remote_algo_path} --name 0 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize 50 --batchTime 500 " >${output_path}0.log &
nohup sshpass ssh -i ${replica2_cert} -n -f ${replica2} ".${remote_algo_path} --name 1 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize 50 --batchTime 500 " >${output_path}1.log &
nohup sshpass ssh -i ${replica3_cert} -n -f ${replica3} ".${remote_algo_path} --name 2 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize 50 --batchTime 500 " >${output_path}2.log &
nohup sshpass ssh -i ${replica4_cert} -n -f ${replica4} ".${remote_algo_path} --name 3 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize 50 --batchTime 500 " >${output_path}3.log &
nohup sshpass ssh -i ${replica5_cert} -n -f ${replica5} ".${remote_algo_path} --name 4 --config ${remote_config_path}  --logFilePath ${remote_log_path} --batchSize 50 --batchTime 500 " >${output_path}4.log &
echo "Started servers"

sleep 20

sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 5 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 1" >${output_path}status1.log
echo "Sent initial status"

sleep 10

echo "Starting client[s]"

nohup sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 5 --defaultReplica 0 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} " >${output_path}5.log &
nohup sshpass ssh -i ${client2_cert} ${client2} ".${remote_ctl_path} --name 6  --defaultReplica 1 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} " >${output_path}6.log &
nohup sshpass ssh -i ${client3_cert} ${client3} ".${remote_ctl_path} --name 7  --defaultReplica 2 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} " >${output_path}7.log &
nohup sshpass ssh -i ${client4_cert} ${client4} ".${remote_ctl_path} --name 8  --defaultReplica 3 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} " >${output_path}8.log &
sshpass ssh -i ${client5_cert} ${client5} ".${remote_ctl_path} --name 9 --defaultReplica 4  --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} " >${output_path}9.log

sleep 10

echo "Completed Client[s]"

sleep 10

sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 5 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 2" >${output_path}status2.log
echo "Sent status to print log"

sleep 50

scp -i ${replica1_cert} ${replica1}:${remote_log_path}0.txt ${output_path}0.txt
scp -i ${replica2_cert} ${replica2}:${remote_log_path}1.txt ${output_path}1.txt
scp -i ${replica3_cert} ${replica3}:${remote_log_path}2.txt ${output_path}2.txt
scp -i ${replica4_cert} ${replica4}:${remote_log_path}3.txt ${output_path}3.txt
scp -i ${replica5_cert} ${replica5}:${remote_log_path}4.txt ${output_path}4.txt

scp -i ${client1_cert} ${client1}:${remote_log_path}5.txt ${output_path}5.txt
scp -i ${client2_cert} ${client2}:${remote_log_path}6.txt ${output_path}6.txt
scp -i ${client3_cert} ${client3}:${remote_log_path}7.txt ${output_path}7.txt
scp -i ${client4_cert} ${client4}:${remote_log_path}8.txt ${output_path}8.txt
scp -i ${client5_cert} ${client5}:${remote_log_path}9.txt ${output_path}9.txt

echo "Copied all the files to local machine"

sleep 5

python3 /home/pasindu/Documents/Raxos/experiments/python/overlay-test.py ${output_path}0.txt ${output_path}1.txt ${output_path}2.txt ${output_path}3.txt ${output_path}4.txt >${output_path}pythonreport.log

dst_directory="/home/pasindu/Desktop/Test/${arrivalRate}/"
mkdir -p "${dst_directory}"
cp -r ${output_path} "${dst_directory}"

sshpass ssh -i ${replica1_cert} ${replica1} "pkill replica"
sshpass ssh -i ${replica2_cert} ${replica2} "pkill replica"
sshpass ssh -i ${replica3_cert} ${replica3} "pkill replica"
sshpass ssh -i ${replica4_cert} ${replica4} "pkill replica"
sshpass ssh -i ${replica5_cert} ${replica5} "pkill replica"

sshpass ssh -i ${client1_cert} ${client1} "pkill client"
sshpass ssh -i ${client2_cert} ${client2} "pkill client"
sshpass ssh -i ${client3_cert} ${client3} "pkill client"
sshpass ssh -i ${client4_cert} ${client4} "pkill client"
sshpass ssh -i ${client5_cert} ${client5} "pkill client"

echo "killed  instances"

echo "Finish test"
