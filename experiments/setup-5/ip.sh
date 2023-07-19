# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem # replace this with your private key
user_name="ubuntu"

replica1_name=ec2-18-176-56-100.ap-northeast-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-13-233-238-185.ap-south-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-18-138-248-160.ap-southeast-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-34-250-227-16.eu-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-18-231-183-28.sa-east-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-52-193-0-185.ap-northeast-1.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-43-205-142-43.ap-south-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-18-141-143-184.ap-southeast-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-3-250-183-221.eu-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-18-228-10-130.sa-east-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
declare -a machines_tune=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1})
echo "ip addresses loaded"

replica1_ip="18.176.56.100"
replica2_ip="13.233.238.185"
replica3_ip="18.138.248.160"
replica4_ip="34.250.227.16"
replica5_ip="18.231.183.28"

client1_ip="52.193.0.185"
client2_ip="43.205.142.43"
client3_ip="18.141.143.184"
client4_ip="3.250.183.221"
client5_ip="18.228.10.130"

device="ens3" #"enp1s0"
