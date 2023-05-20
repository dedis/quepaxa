leaderTimeout=$1
serverMode=$2
leaderMode=$3
pipeline=$4
batchTime=$5
batchSize=$6
arrivalRate=$7
closeLoopWindow=$8
requestPropagationTime=$9

# create a new local configuration
rm -r configuration/local
mkdir -p configuration/local
python3 configuration/config-generate.py 5 5 0.0.0.0 0.0.0.0 0.0.0.0 0.0.0.0 0.0.0.0 0.0.0.0 0.0.0.0 0.0.0.0 0.0.0.0 0.0.0.0 > configuration/local/configuration.yml

# build the project
mage generate && mage build

raxos_path="replica/bin/replica"
ctl_path="client/bin/client"
output_path="logs/${leaderTimeout}/${serverMode}/${leaderMode}/${pipeline}/${batchTime}/${batchSize}/${arrivalRate}/${closeLoopWindow}/${requestPropagationTime}/"

rm -r "${output_path}" ; mkdir -p "${output_path}"

echo "Removed old log files"

pkill replica; pkill replica; pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client; pkill client; pkill client; 

echo "Killed previously running instances"

nohup ./${raxos_path} --name 1 --debugOn --debugLevel 200 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --serverMode "${serverMode}" --requestPropagationTime "${requestPropagationTime}" --logFilePath "${output_path}" >${output_path}1.log &
nohup ./${raxos_path} --name 2 --debugOn --debugLevel 200 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --serverMode "${serverMode}" --requestPropagationTime "${requestPropagationTime}" --logFilePath "${output_path}" >${output_path}2.log &
nohup ./${raxos_path} --name 3 --debugOn --debugLevel 200 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --serverMode "${serverMode}" --requestPropagationTime "${requestPropagationTime}" --logFilePath "${output_path}" >${output_path}3.log &
nohup ./${raxos_path} --name 4 --debugOn --debugLevel 200 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --serverMode "${serverMode}" --requestPropagationTime "${requestPropagationTime}" --logFilePath "${output_path}" >${output_path}4.log &
nohup ./${raxos_path} --name 5 --debugOn --debugLevel 200 --batchSize "${batchSize}" --batchTime  "${batchTime}" --leaderTimeout "${leaderTimeout}" --pipelineLength "${pipeline}" --leaderMode "${leaderMode}" --serverMode "${serverMode}" --requestPropagationTime "${requestPropagationTime}" --logFilePath "${output_path}" >${output_path}5.log &

echo "Started 5 servers"

sleep 3

nohup ./${ctl_path} --name 22 --requestType status --operationType 1  --logFilePath "${output_path}" >${output_path}status1.log &

echo "Sent initial status to bootstrap"

sleep 12

echo "Starting client[s]"

nohup ./${ctl_path} --name 21 --debugOn --debugLevel 100 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${closeLoopWindow}" --logFilePath "${output_path}"  >${output_path}21.log &
nohup ./${ctl_path} --name 22 --debugOn --debugLevel 100 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${closeLoopWindow}" --logFilePath "${output_path}"  >${output_path}22.log &
nohup ./${ctl_path} --name 23 --debugOn --debugLevel 100 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${closeLoopWindow}" --logFilePath "${output_path}"  >${output_path}23.log &
nohup ./${ctl_path} --name 24 --debugOn --debugLevel 100 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${closeLoopWindow}" --logFilePath "${output_path}"  >${output_path}24.log &
nohup ./${ctl_path} --name 25 --debugOn --debugLevel 100 --requestType request --arrivalRate "${arrivalRate}"  --batchSize "${batchSize}" --batchTime "${batchTime}" --window "${closeLoopWindow}" --logFilePath "${output_path}"  >${output_path}25.log &

sleep 100

echo "Completed Client[s]"

nohup ./${ctl_path} --name 23 --requestType status --operationType 4 --logFilePath "${output_path}"  >${output_path}status4.log &
nohup ./${ctl_path} --name 22 --requestType status --operationType 2 --logFilePath "${output_path}" >${output_path}status2.log &

echo "Sent status to print log and print steps per slot"

sleep 40

pkill replica; pkill replica; pkill replica; pkill replica; pkill replica
pkill client; pkill client; pkill client; pkill client; pkill client

rm -r configuration/local

python3 integration-test/python/overlay-test.py ${output_path}1-consensus.txt ${output_path}2-consensus.txt ${output_path}3-consensus.txt ${output_path}4-consensus.txt ${output_path}5-consensus.txt

echo "Killed previously running instances"

echo "Finish test"
