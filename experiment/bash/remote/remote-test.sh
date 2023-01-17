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

rm -r ${output_path}
mkdir ${output_path}

echo "Removed old log files"

sshpass ssh -i ${replica1_cert} ${replica1} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${replica2_cert} ${replica2} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${replica3_cert} ${replica3} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${replica4_cert} ${replica4} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${replica5_cert} ${replica5} "rm ${remote_log_path} && mkdir ${remote_log_path}"

sshpass ssh -i ${client1_cert} ${client1} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${client2_cert} ${client2} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${client3_cert} ${client3} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${client4_cert} ${client4} "rm ${remote_log_path} && mkdir ${remote_log_path}"
sshpass ssh -i ${client5_cert} ${client5} "rm ${remote_log_path} && mkdir ${remote_log_path}"

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

nohup sshpass ssh -i ${replica1_cert} -n -f ${replica1} ".${remote_algo_path} --name 1 --config ${remote_config_path} --batchSize 50 --batchTime 1 --pipelineLength 1 --leaderMode 0 --epochSize 20 --serverMode  0  --logFilePath ${remote_log_path}" >${output_path}1.log &
nohup sshpass ssh -i ${replica2_cert} -n -f ${replica2} ".${remote_algo_path} --name 2 --config ${remote_config_path} --batchSize 50 --batchTime 1 --pipelineLength 1 --leaderMode 0 --epochSize 20 --serverMode  0  --logFilePath ${remote_log_path}" >${output_path}2.log &
nohup sshpass ssh -i ${replica3_cert} -n -f ${replica3} ".${remote_algo_path} --name 3 --config ${remote_config_path} --batchSize 50 --batchTime 1 --pipelineLength 1 --leaderMode 0 --epochSize 20 --serverMode  0  --logFilePath ${remote_log_path}" >${output_path}3.log &
nohup sshpass ssh -i ${replica4_cert} -n -f ${replica4} ".${remote_algo_path} --name 4 --config ${remote_config_path} --batchSize 50 --batchTime 1 --pipelineLength 1 --leaderMode 0 --epochSize 20 --serverMode  0  --logFilePath ${remote_log_path}" >${output_path}4.log &
nohup sshpass ssh -i ${replica5_cert} -n -f ${replica5} ".${remote_algo_path} --name 5 --config ${remote_config_path} --batchSize 50 --batchTime 1 --pipelineLength 1 --leaderMode 0 --epochSize 20 --serverMode  0  --logFilePath ${remote_log_path}" >${output_path}5.log &
echo "Started servers"

sleep 5

sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 1" >${output_path}status1.log
echo "Sent initial status"

sleep 20

echo "Starting client[s]"

nohup sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 50 " >${output_path}21.log &
nohup sshpass ssh -i ${client2_cert} ${client2} ".${remote_ctl_path} --name 22 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 50 " >${output_path}22.log &
nohup sshpass ssh -i ${client3_cert} ${client3} ".${remote_ctl_path} --name 23 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 50 " >${output_path}23.log &
nohup sshpass ssh -i ${client4_cert} ${client4} ".${remote_ctl_path} --name 24 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 50 " >${output_path}24.log &
nohup sshpass ssh -i ${client5_cert} ${client5} ".${remote_ctl_path} --name 25 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType request --arrivalRate ${arrivalRate} --batchSize 50 " >${output_path}25.log &

sleep 200

echo "Completed Client[s]"

sleep 10

sshpass ssh -i ${client1_cert} ${client1} ".${remote_ctl_path} --name 21 --config ${remote_config_path} --logFilePath ${remote_log_path} --requestType status --operationType 2" >${output_path}status2.log
echo "Sent status to print log"

sleep 50

scp -i ${replica1_cert} ${replica1}:${remote_log_path}1-consensus.txt ${output_path}1-consensus.txt
scp -i ${replica1_cert} ${replica1}:${remote_log_path}1-mempool.txt ${output_path}1-mempool.txt

scp -i ${replica2_cert} ${replica2}:${remote_log_path}2-consensus.txt ${output_path}2-consensus.txt
scp -i ${replica2_cert} ${replica2}:${remote_log_path}2-mempool.txt ${output_path}2-mempool.txt

scp -i ${replica3_cert} ${replica3}:${remote_log_path}3-consensus.txt ${output_path}3-consensus.txt
scp -i ${replica3_cert} ${replica3}:${remote_log_path}3-mempool.txt ${output_path}3-mempool.txt

scp -i ${replica4_cert} ${replica4}:${remote_log_path}4-consensus.txt ${output_path}4-consensus.txt
scp -i ${replica4_cert} ${replica4}:${remote_log_path}4-mempool.txt ${output_path}4-mempool.txt

scp -i ${replica5_cert} ${replica5}:${remote_log_path}5-consensus.txt ${output_path}5-consensus.txt
scp -i ${replica5_cert} ${replica5}:${remote_log_path}5-mempool.txt ${output_path}5-mempool.txt

scp -i ${client1_cert} ${client1}:${remote_log_path}21.txt ${output_path}21.txt
scp -i ${client2_cert} ${client2}:${remote_log_path}22.txt ${output_path}22.txt
scp -i ${client3_cert} ${client3}:${remote_log_path}23.txt ${output_path}23.txt
scp -i ${client4_cert} ${client4}:${remote_log_path}24.txt ${output_path}24.txt
scp -i ${client5_cert} ${client5}:${remote_log_path}25.txt ${output_path}25.txt

echo "Copied all the files to local machine"

sleep 5

python3 experiment/python/overlay-test.py ${output_path}1-consensus.txt ${output_path}2-consensus.txt ${output_path}3-consensus.txt ${output_path}4-consensus.txt ${output_path}5-consensus.txt >${output_path}local-test-consensus.log
python3 experiment/python/overlay-test.py ${output_path}1-mempool.txt ${output_path}2-mempool.txt ${output_path}3-mempool.txt ${output_path}4-mempool.txt ${output_path}5-mempool.txt >${output_path}local-test-mempool.log
python3 experiment/python/throughputvstime.py throughput
python3 experiment/python/throughputvstime.py latency

timestamp=$(date +%s)
dst_directory="/home/pasindu/Desktop/Test/${timestamp}/${arrivalRate}/"
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
