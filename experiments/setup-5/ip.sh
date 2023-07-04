# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-184-73-18-9.compute-1.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-54-82-102-216.compute-1.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-52-207-244-195.compute-1.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-3-87-153-249.compute-1.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-163-212-135.compute-1.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-54-89-222-30.compute-1.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-54-157-14-162.compute-1.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-54-91-149-148.compute-1.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-54-164-70-44.compute-1.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-242-255-94.compute-1.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="184.73.18.9"
replica2_ip="54.82.102.216"
replica3_ip="52.207.244.195"
replica4_ip="3.87.153.249"
replica5_ip="54.163.212.135"

client1_ip="54.89.222.30"
client2_ip="54.157.14.162"
client3_ip="54.91.149.148"
client4_ip="54.164.70.44"
client5_ip="54.242.255.94"

device="ens3" #"enp1s0"
