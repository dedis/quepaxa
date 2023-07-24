# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-35-93-7-223.us-west-2.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-35-165-6-8.us-west-2.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-54-184-246-145.us-west-2.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-54-212-114-230.us-west-2.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-34-222-4-7.us-west-2.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-54-212-103-170.us-west-2.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-54-212-103-170.us-west-2.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-54-212-103-170.us-west-2.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-54-212-103-170.us-west-2.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-212-103-170.us-west-2.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
declare -a machines_tune=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1})
echo "ip addresses loaded"

replica1_ip="35.93.7.223"
replica2_ip="35.165.6.8"
replica3_ip="54.184.246.145"
replica4_ip="54.212.114.230"
replica5_ip="34.222.4.7"

client1_ip="54.212.103.170"
client2_ip="3.110.177.243"
client3_ip="13.229.139.178"
client4_ip="3.253.155.13"
client5_ip="52.67.23.103"

device="ens3" #"enp1s0"
