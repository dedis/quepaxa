# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-18-179-36-33.ap-northeast-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-13-232-0-159.ap-south-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-52-221-179-245.ap-southeast-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-3-253-152-110.eu-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-54-233-200-63.sa-east-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-13-231-195-120.ap-northeast-1.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-3-109-60-50.ap-south-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-13-229-224-7.ap-southeast-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-52-50-198-188.eu-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-18-229-158-236.sa-east-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="18.179.36.33"
replica2_ip="13.232.0.159"
replica3_ip="52.221.179.245"
replica4_ip="3.253.152.110"
replica5_ip="54.233.200.63"

client1_ip="13.231.195.120"
client2_ip="3.109.60.50"
client3_ip="13.229.224.7"
client4_ip="52.50.198.188"
client5_ip="18.229.158.236"

device="ens3" #"enp1s0"
