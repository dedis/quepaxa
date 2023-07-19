# this file assigns the ip address of each replica and client to variables
pwd=$(pwd)
cert=${pwd}/experiments/private_key_aws/pasindu_rsa
user_name="ubuntu"

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

replica6=pasindu@dedis-140.icsil1.epfl.ch
replica6_name=dedis-140.icsil1.epfl.ch

replica7=pasindu@dedis-141.icsil1.epfl.ch
replica7_name=dedis-141.icsil1.epfl.ch

replica8=pasindu@dedis-142.icsil1.epfl.ch
replica8_name=dedis-142.icsil1.epfl.ch

replica9=pasindu@dedis-143.icsil1.epfl.ch
replica9_name=dedis-143.icsil1.epfl.ch

replica10=pasindu@dedis-144.icsil1.epfl.ch
replica10_name=dedis-144.icsil1.epfl.ch

replica11=pasindu@dedis-142.icsil1.epfl.ch
replica11_name=dedis-142.icsil1.epfl.ch

replica12=pasindu@dedis-143.icsil1.epfl.ch
replica12_name=dedis-143.icsil1.epfl.ch

replica13=pasindu@dedis-144.icsil1.epfl.ch
replica13_name=dedis-144.icsil1.epfl.ch


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

client6=pasindu@dedis-145.icsil1.epfl.ch
client6_name=dedis-145.icsil1.epfl.ch

client7=pasindu@dedis-146.icsil1.epfl.ch
client7_name=dedis-146.icsil1.epfl.ch

client8=pasindu@dedis-147.icsil1.epfl.ch
client8_name=dedis-147.icsil1.epfl.ch

client9=pasindu@dedis-148.icsil1.epfl.ch
client9_name=dedis-148.icsil1.epfl.ch

client10=pasindu@dedis-149.icsil1.epfl.ch
client10_name=dedis-149.icsil1.epfl.ch

client11=pasindu@dedis-145.icsil1.epfl.ch
client11_name=dedis-145.icsil1.epfl.ch

client12=pasindu@dedis-146.icsil1.epfl.ch
client12_name=dedis-146.icsil1.epfl.ch

client13=pasindu@dedis-147.icsil1.epfl.ch
client13_name=dedis-147.icsil1.epfl.ch

declare -a machines=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${replica6} ${replica7} ${replica8} ${replica9} ${replica10}  ${replica11} ${replica12} ${replica13} ${client1} ${client2} ${client3} ${client4} ${client5} ${client6} ${client7} ${client8} ${client9} ${client10} ${client11} ${client12} ${client13})
echo "ip addresses loaded"

replica1_ip="10.156.33.140"
replica2_ip="10.156.33.141"
replica3_ip="10.156.33.142"
replica4_ip="10.156.33.143"
replica5_ip="10.156.33.144"
replica6_ip="10.156.33.140"
replica7_ip="10.156.33.141"
replica8_ip="10.156.33.142"
replica9_ip="10.156.33.143"
replica10_ip="10.156.33.144"
replica11_ip="10.156.33.140"
replica12_ip="10.156.33.141"
replica13_ip="10.156.33.142"

client1_ip="10.156.33.145"
client2_ip="10.156.33.146"
client3_ip="10.156.33.147"
client4_ip="10.156.33.148"
client5_ip="10.156.33.149"
client6_ip="10.156.33.145"
client7_ip="10.156.33.146"
client8_ip="10.156.33.147"
client9_ip="10.156.33.148"
client10_ip="10.156.33.149"
client11_ip="10.156.33.145"
client12_ip="10.156.33.146"
client13_ip="10.156.33.147"