# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-18-234-52-145.compute-1.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-52-90-99-154.compute-1.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-52-91-91-116.compute-1.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-23-20-216-248.compute-1.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-107-20-83-93.compute-1.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-34-202-163-199.compute-1.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-18-215-144-141.compute-1.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-54-226-177-182.compute-1.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-54-205-220-226.compute-1.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=eec2-54-210-252-84.compute-1.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
declare -a machines_tune=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1})
echo "ip addresses loaded"

replica1_ip="18.234.52.145"
replica2_ip="52.90.99.154"
replica3_ip="52.91.91.116"
replica4_ip="23.20.216.248"
replica5_ip="107.20.83.93"

client1_ip="34.202.163.199"
client2_ip="18.215.144.141"
client3_ip="54.226.177.182"
client4_ip="54.205.220.226"
client5_ip="54.210.252.84"

device="ens3" #"enp1s0"
