pwd=$(pwd)
. "${pwd}"/experiments/setup-13/ip.sh

reset_directory="rm -r /home/${user_name}/*; mkdir -p /home/${user_name}/raxos; mkdir -p /home/${user_name}/raxos/logs"
local_binary_path="${pwd}/experiments/binary/"
replica_home_path="/home/${user_name}/raxos/"

# generate config files

python3 experiments/python/quepaxa-config.py 13 13 ${replica1_ip} ${replica2_ip} ${replica3_ip} ${replica4_ip} ${replica5_ip} ${replica6_ip} ${replica7_ip} ${replica8_ip} ${replica9_ip} ${replica10_ip} ${replica11_ip} ${replica12_ip} ${replica13_ip} ${client1_ip} ${client2_ip} ${client3_ip} ${client4_ip} ${client5_ip} ${client6_ip} ${client7_ip} ${client8_ip} ${client9_ip} ${client10_ip} ${client11_ip} ${client12_ip} ${client13_ip} > experiments/binary/quepaxa.yml



for i in "${machines[@]}"
do
   echo "copying files to ${i}"
   sshpass ssh "$i" -i ${cert} "${reset_directory}"
   scp -r -i ${cert} ${local_binary_path} "$i":${replica_home_path}
done


rm experiments/binary/quepaxa.yml

echo "13 replica setup complete"

