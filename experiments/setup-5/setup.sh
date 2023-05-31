pwd=$(pwd)
reset_directory="rm -r /home/ubuntu/raxos; mkdir /home/ubuntu/raxos; mkdir -p /home/ubuntu/raxos/logs"
local_zip_path="${pwd}/bin/"
replica_home_path="/home/ubuntu/raxos/"
raxos_config="${pwd}/aws/best-case-lan/raxos/configuration/raxos_configuration.yml"
paxos_raft_config="${pwd}/aws/best-case-lan/paxos_raft/configuration/paxos_raft_configuration.yml"


. "${pwd}"/aws/best-case-lan/setup/ip.sh

for i in "${machines[@]}"
do
   echo "$i"
   sshpass ssh "$i" -i ${cert} "${reset_directory}"
   scp -r -i ${cert} ${local_zip_path} "$i":${replica_home_path}
   scp -i ${cert} ${raxos_config} "$i":${replica_home_path}
   scp -i ${cert} ${paxos_raft_config} "$i":${replica_home_path}
done


