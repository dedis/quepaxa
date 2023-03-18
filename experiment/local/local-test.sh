epoch=$1
leaderTimeout=$2
serverMode=$3
leaderMode=$4
pipeline=$5
batchTime=$6
batchSize=$7
arrivalRate=$8

mage generate && mage build

raxos_path="replica/bin/replica"
ctl_path="client/bin/client"
output_path="logs/"

rm -r ${output_path} ; mkdir ${output_path}

echo "Removed old log files"

pkill replica; pkill replica; pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client; pkill client; pkill client; 

echo "Killed previously running instances"

nohup ./${raxos_path} --name 1 --debugOn --debugLevel 30 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --epochSize "${epoch}" --serverMode "${serverMode}" >${output_path}1.log &
nohup ./${raxos_path} --name 2 --debugOn --debugLevel 30 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --epochSize "${epoch}" --serverMode "${serverMode}" >${output_path}2.log &
nohup ./${raxos_path} --name 3 --debugOn --debugLevel 30 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --epochSize "${epoch}" --serverMode "${serverMode}" >${output_path}3.log &
nohup ./${raxos_path} --name 4 --debugOn --debugLevel 30 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --epochSize "${epoch}" --serverMode "${serverMode}" >${output_path}4.log &
nohup ./${raxos_path} --name 5 --debugOn --debugLevel 30 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --epochSize "${epoch}" --serverMode "${serverMode}" >${output_path}5.log &

echo "Started 5 servers"

sleep 3

nohup ./${ctl_path} --name 22 --requestType status --operationType 1 >${output_path}status1.log &

echo "Sent initial status to bootstrap"

sleep 3

echo "Starting client[s]"

nohup ./${ctl_path} --name 21 --debugOn --debugLevel 20 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" >${output_path}21.log &


sleep 200

echo "Completed Client[s]"

nohup ./${ctl_path} --name 22 --requestType status --operationType 2 >${output_path}status2.log &

echo "Sent status to print log"

sleep 10

pkill replica; pkill replica; pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client; pkill client; pkill client

python3 experiment/python/overlay-test.py 60 ${output_path}1-consensus.txt ${output_path}2-consensus.txt ${output_path}3-consensus.txt ${output_path}4-consensus.txt ${output_path}5-consensus.txt >${output_path}local-test-consensus.log
python3 experiment/python/overlay-test.py 60 ${output_path}1-mempool.txt   ${output_path}2-mempool.txt   ${output_path}3-mempool.txt   ${output_path}4-mempool.txt   ${output_path}5-mempool.txt   >${output_path}local-test-mempool.log

python3 experiment/python/throughputvstime.py throughput
python3 experiment/python/throughputvstime.py latency

echo "Killed previously running instances"

mkdir -p    /home/pasindu/Desktop/raxos/"${epoch}"/"${leaderTimeout}"/"${serverMode}"/"${leaderMode}"/"${pipeline}"/"${batchTime}"/"${batchSize}"/"${arrivalRate}"/logs/
cp -r logs/ /home/pasindu/Desktop/raxos/"${epoch}"/"${leaderTimeout}"/"${serverMode}"/"${leaderMode}"/"${pipeline}"/"${batchTime}"/"${batchSize}"/"${arrivalRate}"/logs/

echo "Finish test"
