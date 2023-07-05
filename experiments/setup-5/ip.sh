# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1_name=ec2-43-207-88-116.ap-northeast-1.compute.amazonaws.com
replica1=${user_name}@${replica1_name}

replica2_name=ec2-13-126-97-245.ap-south-1.compute.amazonaws.com
replica2=${user_name}@${replica2_name}

replica3_name=ec2-3-0-55-141.ap-southeast-1.compute.amazonaws.com
replica3=${user_name}@${replica3_name}

replica4_name=ec2-34-247-217-185.eu-west-1.compute.amazonaws.com
replica4=${user_name}@${replica4_name}

replica5_name=ec2-18-231-115-217.sa-east-1.compute.amazonaws.com
replica5=${user_name}@${replica5_name}

client1_name=ec2-43-207-204-90.ap-northeast-1.compute.amazonaws.com
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
echo "ip addresses loaded"

replica1_ip="43.207.88.116"
replica2_ip="13.126.97.245"
replica3_ip="3.0.55.141"
replica4_ip="34.247.217.185"
replica5_ip="18.231.115.217"

client1_ip="43.207.204.90"
client2_ip="13.233.148.223"
client3_ip="18.140.58.201"
client4_ip="3.249.136.78"
client5_ip="54.233.98.87"

device="ens3" #"enp1s0"
