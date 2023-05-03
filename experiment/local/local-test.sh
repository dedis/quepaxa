epoch=$1
leaderTimeout=$2
serverMode=$3
leaderMode=$4
pipeline=$5
batchTime=$6
batchSize=$7
arrivalRate=$8
window=$9

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

sleep 15

echo "Starting client[s]"

nohup ./${ctl_path} --name 21 --debugOn --debugLevel 20 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${window}" >${output_path}21.log &
nohup ./${ctl_path} --name 22 --debugOn --debugLevel 20 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${window}" >${output_path}22.log &
nohup ./${ctl_path} --name 23 --debugOn --debugLevel 20 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${window}" >${output_path}23.log &
nohup ./${ctl_path} --name 24 --debugOn --debugLevel 20 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${window}" >${output_path}24.log &
nohup ./${ctl_path} --name 25 --debugOn --debugLevel 20 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${window}" >${output_path}25.log &

sleep 80

echo "Completed Client[s]"

nohup ./${ctl_path} --name 23 --requestType status --operationType 4 >${output_path}status4.log &
nohup ./${ctl_path} --name 22 --requestType status --operationType 2 >${output_path}status2.log &

echo "Sent status to print log and print steps per slot"

sleep 40

python3 experiment/python/overlay-test.py 60 ${output_path}1-consensus.txt ${output_path}2-consensus.txt ${output_path}3-consensus.txt ${output_path}4-consensus.txt ${output_path}5-consensus.txt >${output_path}local-test-consensus.log
python3 experiment/python/overlay-test.py 60 ${output_path}1-mempool.txt   ${output_path}2-mempool.txt   ${output_path}3-mempool.txt   ${output_path}4-mempool.txt   ${output_path}5-mempool.txt   >${output_path}local-test-mempool.log

pkill replica; pkill replica; pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client; pkill client; pkill client

echo "Killed previously running instances"

echo "Finish test"
