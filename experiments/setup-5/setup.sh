pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

reset_directory="rm -r /home/${user_name}/*; mkdir -p /home/${user_name}/raxos; mkdir -p /home/${user_name}/raxos/logs"
local_binary_path="${pwd}/experiments/binary/"
replica_home_path="/home/${user_name}/raxos/"

# generate config files

python3 experiments/python/quepaxa-config.py 5 5 ${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5} > experiments/binary/quepaxa.yml
python3 experiments/python/paxos-config.py   5 5 ${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5} > experiments/binary/paxos.yml


for i in "${machines[@]}"
do
   echo "copying files to ${i}"
   sshpass ssh "$i" -i ${cert} "${reset_directory}"
   scp -r -i ${cert} ${local_binary_path} "$i":${replica_home_path}
done


rm experiments/binary/quepaxa.yml
rm experiments/binary/paxos.yml

echo "5 replica setup complete"

