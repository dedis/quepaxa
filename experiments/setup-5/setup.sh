pwd=$(pwd)
. "${pwd}"/experiments/setup-5/ip.sh

reset_directory="rm -r /home/${user_name}/*; mkdir -p /home/${user_name}/raxos; mkdir -p /home/${user_name}/raxos/logs"
local_binary_path="${pwd}/experiments/binary/"
replica_home_path="/home/${user_name}/raxos/"

# generate config files

python3 experiments/python/quepaxa-config.py 5 5 ${replica1_ip} ${replica2_ip} ${replica3_ip} ${replica4_ip} ${replica5_ip} ${client1_ip} ${client2_ip} ${client3_ip} ${client4_ip} ${client5_ip} > experiments/binary/quepaxa.yml
python3 experiments/python/paxos-config.py   5 5 ${replica1_ip} ${replica2_ip} ${replica3_ip} ${replica4_ip} ${replica5_ip} ${client1_ip} ${client2_ip} ${client3_ip} ${client4_ip} ${client5_ip} > experiments/binary/paxos.yml


for i in "${machines[@]}"
do
   echo "copying files to ${i}"
   sshpass ssh "$i" -i ${cert} "${reset_directory}"
   scp -r -i ${cert} ${local_binary_path} "$i":${replica_home_path}
done


rm experiments/binary/quepaxa.yml
rm experiments/binary/paxos.yml

echo "5 replica setup complete"

