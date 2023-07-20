# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem # replace this with your private key
user_name="ubuntu"

replica1_name=ec2-54-250-204-196.ap-northeast-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-3-109-5-69.ap-south-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-18-140-51-67.ap-southeast-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-3-253-252-196.eu-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-18-228-190-178.sa-east-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-52-68-160-14.ap-northeast-1.compute.amazonaws.com
client1=${user_name}@${client1_name}

client2_name=ec2-13-232-6-64.ap-south-1.compute.amazonaws.com
client2=${user_name}@${client2_name}

client3_name=ec2-13-213-73-50.ap-southeast-1.compute.amazonaws.com
client3=${user_name}@${client3_name}

client4_name=ec2-34-254-159-194.eu-west-1.compute.amazonaws.com
client4=${user_name}@${client4_name}

client5_name=ec2-18-229-163-11.sa-east-1.compute.amazonaws.com
client5=${user_name}@${client5_name}

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
declare -a machines_tune=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1})
echo "ip addresses loaded"

replica1_ip="54.250.204.196"
replica2_ip="3.109.5.69"
replica3_ip="18.140.51.67"
replica4_ip="3.253.252.196"
replica5_ip="18.228.190.178"

client1_ip="52.68.160.14"
client2_ip="13.232.6.64"
client3_ip="13.213.73.50"
client4_ip="34.254.159.194"
client5_ip="18.229.163.11"

device="ens3" #"enp1s0"
