# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-35-92-243-9.us-west-2.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-34-215-234-137.us-west-2.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-54-201-143-251.us-west-2.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-54-69-99-239.us-west-2.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-35-92-219-76.us-west-2.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-54-212-125-133.us-west-2.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-13-233-148-223.ap-south-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-18-140-58-201.ap-southeast-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-3-249-136-78.eu-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-54-233-98-87.sa-east-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
declare -a machines_tune=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1})
echo "ip addresses loaded"

replica1_ip="35.92.243.9"
replica2_ip="34.215.234.137"
replica3_ip="54.201.143.251"
replica4_ip="54.69.99.239"
replica5_ip="35.92.219.76"

client1_ip="54.212.125.133"
client2_ip="54.212.125.133"
client3_ip="18.140.58.201"
client4_ip="3.249.136.78"
client5_ip="54.233.98.87"

device="ens3" #"enp1s0"
