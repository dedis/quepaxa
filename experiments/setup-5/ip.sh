# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/private.pem
user_name="ubuntu"

replica1=${user_name}@ec2-54-172-197-20.compute-1.amazonaws.com
replica1_name=ec2-54-172-197-20.compute-1.amazonaws.com

replica2=${user_name}@ec2-54-234-146-41.compute-1.amazonaws.com
replica2_name=ec2-54-234-146-41.compute-1.amazonaws.com

replica3=${user_name}@ec2-54-234-241-51.compute-1.amazonaws.com
replica3_name=ec2-54-234-241-51.compute-1.amazonaws.com

replica4=${user_name}@ec2-3-95-67-49.compute-1.amazonaws.com
replica4_name=ec2-3-95-67-49.compute-1.amazonaws.com

replica5=${user_name}@ec2-54-87-63-253.compute-1.amazonaws.com
replica5_name=ec2-54-87-63-253.compute-1.amazonaws.com

client1=${user_name}@ec2-52-23-201-149.compute-1.amazonaws.com
client1_name=ec2-52-23-201-149.compute-1.amazonaws.com

client2=${user_name}@ec2-54-146-146-149.compute-1.amazonaws.com
client2_name=ec2-54-146-146-149.compute-1.amazonaws.com

client3=${user_name}@ec2-34-234-96-45.compute-1.amazonaws.com
client3_name=ec2-34-234-96-45.compute-1.amazonaws.com

client4=${user_name}@ec2-54-173-187-230.compute-1.amazonaws.com
client4_name=ec2-54-173-187-230.compute-1.amazonaws.com

client5=${user_name}@ec2-54-91-179-43.compute-1.amazonaws.com
client5_name=ec2-54-91-179-43.compute-1.amazonaws.com

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="54.172.197.20"
replica2_ip="54.234.146.41"
replica3_ip="54.234.241.51"
replica4_ip="3.95.67.49"
replica5_ip="54.87.63.253"

client1_ip="52.23.201.149"
client2_ip="54.146.146.149"
client3_ip="34.234.96.45"
client4_ip="54.173.187.230"
client5_ip="54.91.179.43"

device="ens3" #"enp1s0"
