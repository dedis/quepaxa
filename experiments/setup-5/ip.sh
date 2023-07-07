# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-54-168-49-240.ap-northeast-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-13-127-107-189.ap-south-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-54-179-179-67.ap-southeast-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-54-74-164-158.eu-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-15-228-220-219.sa-east-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-52-195-168-183.ap-northeast-1.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-13-232-82-187.ap-south-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-18-139-224-187.ap-southeast-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-18-203-67-51.eu-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-18-231-109-31.sa-east-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
declare -a machines_tune=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1})
echo "ip addresses loaded"

replica1_ip="54.168.49.240"
replica2_ip="13.127.107.189"
replica3_ip="54.179.179.67"
replica4_ip="54.74.164.158"
replica5_ip="15.228.220.219"

client1_ip="52.195.168.183"
client2_ip="13.232.82.187"
client3_ip="18.139.224.187"
client4_ip="18.203.67.51"
client5_ip="18.231.109.31"

device="ens3" #"enp1s0"
