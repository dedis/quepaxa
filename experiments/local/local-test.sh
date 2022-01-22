mage generate && mage build

raxos_path="replica/bin/replica"
ctl_path="client/bin/client"
output_path="logs/"

rm nohup.out
rm ${output_path}0.txt
rm ${output_path}1.txt
rm ${output_path}2.txt
rm ${output_path}3.txt
rm ${output_path}4.txt
rm ${output_path}5.txt

pkill replica
pkill replica
pkill replica
pkill client
pkill client
pkill client

nohup ./${raxos_path} --name 0 &
nohup ./${raxos_path} --name 1 &
nohup ./${raxos_path} --name 2 &

echo "Started servers, Please check the nohup.out"

sleep 10

./${ctl_path}  --name 3 --requestType status --operationType 1

echo "Sent initial status"

echo "Starting client[s]"

nohup ./${ctl_path}  --name 3 --defaultReplica 0 --requestType request &
nohup ./${ctl_path}  --name 4 --defaultReplica 1 --requestType request &
nohup ./${ctl_path}  --name 5 --defaultReplica 2 --requestType request &

sleep 150

echo "Completed Client[s]"

./${ctl_path}  --name 3 --requestType status --operationType 2

echo "Sent status to print log"

sleep 10

pkill replica
pkill replica
pkill replica
pkill client
pkill client
pkill client

echo "Finish test"