# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-18-176-215-125.ap-northeast-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-3-109-183-37.ap-south-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-3-0-51-53.ap-southeast-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-54-216-10-235.eu-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-18-228-219-145.sa-east-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-43-207-237-241.ap-northeast-1.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-3-110-177-243.ap-south-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-13-229-139-178.ap-southeast-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-3-253-155-13.eu-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-52-67-23-103.sa-east-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
declare -a machines_tune=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1})
echo "ip addresses loaded"

replica1_ip="18.176.215.125"
replica2_ip="3.109.183.37"
replica3_ip="3.0.51.53"
replica4_ip="54.216.10.235"
replica5_ip="18.228.219.145"

client1_ip="43.207.237.241"
client2_ip="3.110.177.243"
client3_ip="13.229.139.178"
client4_ip="3.253.155.13"
client5_ip="52.67.23.103"

device="ens3" #"enp1s0"
