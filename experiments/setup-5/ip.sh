# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu2023.pem
user_name="ubuntu"

replica1=${user_name}@ec2-3-112-210-151.ap-northeast-1.compute.amazonaws.com
replica1_name=ec2-3-112-210-151.ap-northeast-1.compute.amazonaws.com

replica2=${user_name}@ec2-3-109-144-150.ap-south-1.compute.amazonaws.com
replica2_name=ec2-3-109-144-150.ap-south-1.compute.amazonaws.com

replica3=${user_name}@ec2-54-169-115-28.ap-southeast-1.compute.amazonaws.com
replica3_name=ec2-54-169-115-28.ap-southeast-1.compute.amazonaws.com

replica4=${user_name}@ec2-52-50-51-40.eu-west-1.compute.amazonaws.com
replica4_name=ec2-52-50-51-40.eu-west-1.compute.amazonaws.com

replica5=${user_name}@ec2-54-94-32-33.sa-east-1.compute.amazonaws.com
replica5_name=ec2-54-94-32-33.sa-east-1.compute.amazonaws.com

client1=${user_name}@ec2-54-249-7-175.ap-northeast-1.compute.amazonaws.com
client1_name=ec2-54-249-7-175.ap-northeast-1.compute.amazonaws.com

client2=${user_name}@ec2-13-233-88-111.ap-south-1.compute.amazonaws.com
client2_name=ec2-13-233-88-111.ap-south-1.compute.amazonaws.com

client3=${user_name}@ec2-13-215-153-159.ap-southeast-1.compute.amazonaws.com
client3_name=ec2-13-215-153-159.ap-southeast-1.compute.amazonaws.com

client4=${user_name}@ec2-63-32-89-0.eu-west-1.compute.amazonaws.com
client4_name=ec2-63-32-89-0.eu-west-1.compute.amazonaws.com

client5=${user_name}@ec2-18-231-178-141.sa-east-1.compute.amazonaws.com
client5_name=ec2-18-231-178-141.sa-east-1.compute.amazonaws.com

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="3.112.210.151"
replica2_ip="3.109.144.150"
replica3_ip="54.169.115.28"
replica4_ip="52.50.51.40"
replica5_ip="54.94.32.33"

client1_ip="54.249.7.175"
client2_ip="13.233.88.111"
client3_ip="13.215.153.159"
client4_ip="63.32.89.0"
client5_ip="18.231.178.141"

device="ens3" #"enp1s0"
