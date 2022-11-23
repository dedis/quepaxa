arrivalRate=$1
# A local test that
#     1. Build the project
#     2. Spawn 5 replicas
#     3. Boot strap servers
#     4. Spawn a client
#     5. Compare the logs
#     6. Kill instances and clients

mage generate && mage build

raxos_path="replica/bin/replica"
ctl_path="client/bin/client"
output_path="logs/"

rm ${output_path}1.log
rm ${output_path}2.log
rm ${output_path}3.log
rm ${output_path}4.log
rm ${output_path}5.log


rm ${output_path}1-consensus.txt
rm ${output_path}2-consensus.txt
rm ${output_path}3-consensus.txt
rm ${output_path}4-consensus.txt
rm ${output_path}5-consensus.txt

rm ${output_path}1-mempool.txt
rm ${output_path}2-mempool.txt
rm ${output_path}3-mempool.txt
rm ${output_path}4-mempool.txt
rm ${output_path}5-mempool.txt

rm ${output_path}21.txt
rm ${output_path}21.log

rm ${output_path}local-test-consensus.log
rm ${output_path}local-test-mempool.log

rm ${output_path}status1.log
rm ${output_path}status2.log

echo "Removed old log files"

pkill replica
pkill replica
pkill replica
pkill replica
pkill replica
pkill client

echo "Killed previously running instances"

nohup ./${raxos_path} --name 1 --debugOn --debugLevel 8 --batchSize 50 --pipelineLength 1 --leaderMode 0 >${output_path}1.log &
nohup ./${raxos_path} --name 2 --debugOn --debugLevel 8 --batchSize 50 --pipelineLength 1 --leaderMode 0 >${output_path}2.log &
nohup ./${raxos_path} --name 3 --debugOn --debugLevel 8 --batchSize 50 --pipelineLength 1 --leaderMode 0 >${output_path}3.log &
nohup ./${raxos_path} --name 4 --debugOn --debugLevel 8 --batchSize 50 --pipelineLength 1 --leaderMode 0 >${output_path}4.log &
nohup ./${raxos_path} --name 5 --debugOn --debugLevel 8 --batchSize 50 --pipelineLength 1 --leaderMode 0 >${output_path}5.log &

echo "Started 5 servers"

sleep 3

./${ctl_path} --name 21 --requestType status --operationType 1 >${output_path}status1.log

echo "Sent initial status to bootstrap"

sleep 3

echo "Starting client[s]"

nohup ./${ctl_path} --name 21 --debugOn --debugLevel 4 --requestType request --arrivalRate "${arrivalRate}" --batchSize 50 >${output_path}21.log &

sleep 200

echo "Completed Client[s]"

./${ctl_path} --name 21 --requestType status --operationType 2 >${output_path}status2.log

echo "Sent status to print log"

sleep 50

python3 experiment/python/overlay-test.py ${output_path}5-consensus.txt ${output_path}1-consensus.txt ${output_path}2-consensus.txt ${output_path}3-consensus.txt ${output_path}4-consensus.txt >${output_path}local-test-consensus.log
python3 experiment/python/overlay-test.py ${output_path}5-mempool.txt ${output_path}1-mempool.txt ${output_path}2-mempool.txt ${output_path}3-mempool.txt ${output_path}4-mempool.txt >${output_path}local-test-mempool.log

pkill replica
pkill replica
pkill replica
pkill replica
pkill replica
pkill client

echo "Killed previously running instances"

echo "Finish test"
