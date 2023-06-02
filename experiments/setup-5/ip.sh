# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu_rsa

replica1=pasindu@dedis-140.icsil1.epfl.ch
replica1_name=dedis-140.icsil1.epfl.ch

replica2=pasindu@dedis-141.icsil1.epfl.ch
replica2_name=dedis-141.icsil1.epfl.ch

replica3=pasindu@dedis-142.icsil1.epfl.ch
replica3_name=dedis-142.icsil1.epfl.ch

replica4=pasindu@dedis-143.icsil1.epfl.ch
replica4_name=dedis-143.icsil1.epfl.ch

replica5=pasindu@dedis-144.icsil1.epfl.ch
replica5_name=dedis-144.icsil1.epfl.ch

client1=pasindu@dedis-145.icsil1.epfl.ch
client1_name=dedis-145.icsil1.epfl.ch

client2=pasindu@dedis-146.icsil1.epfl.ch
client2_name=dedis-146.icsil1.epfl.ch

client3=pasindu@dedis-147.icsil1.epfl.ch
client3_name=dedis-147.icsil1.epfl.ch

client4=pasindu@dedis-148.icsil1.epfl.ch
client4_name=dedis-148.icsil1.epfl.ch

client5=pasindu@dedis-149.icsil1.epfl.ch
client5_name=dedis-149.icsil1.epfl.ch

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${client1} ${client2} ${client3} ${client4} ${client5})
echo "ip addresses loaded"

replica1_ip="10.156.33.140"
replica2_ip="10.156.33.141"
replica3_ip="10.156.33.142"
replica4_ip="10.156.33.143"
replica5_ip="10.156.33.144"

client1_ip="10.156.33.145"
client2_ip="10.156.33.146"
client3_ip="10.156.33.147"
client4_ip="10.156.33.148"
client5_ip="10.156.33.149"

user_name="pasindu"
device="enp1s0" # "ens3"
