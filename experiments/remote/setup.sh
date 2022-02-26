#     1. Build the project
#     2. Copy client/ configuration/ experiments/ logs/ replica/ directories to logs/raxos/ directory
#     2. Zip logs/raxos/ and create raxos.zip
#     3. Copy raxos.zip to all machines

rm client/bin/client
rm replica/bin/replica

echo "Removed old binaries"
mage generate && mage build
echo "Built Raxos"

reset_directory="rm -r /home/pasindu/raxos; mkdir /home/pasindu/raxos"
kill_insstances="pkill replica ; pkill client"
unzip_raxos="cd /home/pasindu/raxos && unzip raxos.zip"

replica1=pasindu@dedis-140.icsil1.epfl.ch
replica1_cert="/home/pasindu/Pictures/pasindu_rsa"
replica2=pasindu@dedis-141.icsil1.epfl.ch
replica2_cert="/home/pasindu/Pictures/pasindu_rsa"
replica3=pasindu@dedis-142.icsil1.epfl.ch
replica3_cert="/home/pasindu/Pictures/pasindu_rsa"
replica4=pasindu@dedis-143.icsil1.epfl.ch
replica4_cert="/home/pasindu/Pictures/pasindu_rsa"
replica5=pasindu@dedis-144.icsil1.epfl.ch
replica5_cert="/home/pasindu/Pictures/pasindu_rsa"

client1=pasindu@dedis-145.icsil1.epfl.ch
client1_cert="/home/pasindu/Pictures/pasindu_rsa"
client2=pasindu@dedis-146.icsil1.epfl.ch
client2_cert="/home/pasindu/Pictures/pasindu_rsa"
client3=pasindu@dedis-147.icsil1.epfl.ch
client3_cert="/home/pasindu/Pictures/pasindu_rsa"
client4=pasindu@dedis-148.icsil1.epfl.ch
client4_cert="/home/pasindu/Pictures/pasindu_rsa"
client5=pasindu@dedis-149.icsil1.epfl.ch
client5_cert="/home/pasindu/Pictures/pasindu_rsa"

mkdir temp

cp -r logs/ temp/
cp -r client/ temp/
cp -r configuration/ temp/
cp -r experiments/ temp/
cp -r replica/ temp/

zip -r logs/raxos.zip temp/
rm -r temp/

local_zip_path="logs/raxos.zip"
replica_home_path="/home/pasindu/raxos/"

echo "Replica 1"
sshpass ssh ${replica1} -i ${replica1_cert} ${reset_directory}
sshpass ssh ${replica1} -i ${replica1_cert} ${kill_insstances}
scp -i ${replica1_cert} ${local_zip_path} ${replica1}:${replica_home_path}
sshpass ssh ${replica1} -i ${replica1_cert} ${unzip_raxos}

echo "Replica 2"
sshpass ssh ${replica2} -i ${replica2_cert} ${reset_directory}
sshpass ssh ${replica2} -i ${replica2_cert} ${kill_insstances}
scp -i ${replica2_cert} ${local_zip_path} ${replica2}:${replica_home_path}
sshpass ssh ${replica2} -i ${replica2_cert} ${unzip_raxos}

echo "Replica 3"
sshpass ssh ${replica3} -i ${replica3_cert} ${reset_directory}
sshpass ssh ${replica3} -i ${replica3_cert} ${kill_insstances}
scp -i ${replica3_cert} ${local_zip_path} ${replica3}:${replica_home_path}
sshpass ssh ${replica3} -i ${replica3_cert} ${unzip_raxos}

echo "Replica 4"
sshpass ssh ${replica4} -i ${replica4_cert} ${reset_directory}
sshpass ssh ${replica4} -i ${replica4_cert} ${kill_insstances}
scp -i ${replica4_cert} ${local_zip_path} ${replica4}:${replica_home_path}
sshpass ssh ${replica4} -i ${replica4_cert} ${unzip_raxos}

echo "Replica 5"
sshpass ssh ${replica5} -i ${replica5_cert} ${reset_directory}
sshpass ssh ${replica5} -i ${replica5_cert} ${kill_insstances}
scp -i ${replica5_cert} ${local_zip_path} ${replica5}:${replica_home_path}
sshpass ssh ${replica5} -i ${replica5_cert} ${unzip_raxos}

echo "Client 1"
sshpass ssh ${client1} -i ${client1_cert} ${reset_directory}
sshpass ssh ${client1} -i ${client1_cert} ${kill_insstances}
scp -i ${client1_cert} ${local_zip_path} ${client1}:${replica_home_path}
sshpass ssh ${client1} -i ${client1_cert} ${unzip_raxos}

echo "Client 2"
sshpass ssh ${client2} -i ${client2_cert} ${reset_directory}
sshpass ssh ${client2} -i ${client2_cert} ${kill_insstances}
scp -i ${client2_cert} ${local_zip_path} ${client2}:${replica_home_path}
sshpass ssh ${client2} -i ${client2_cert} ${unzip_raxos}

echo "Client 3"
sshpass ssh ${client3} -i ${client3_cert} ${reset_directory}
sshpass ssh ${client3} -i ${client3_cert} ${kill_insstances}
scp -i ${client3_cert} ${local_zip_path} ${client3}:${replica_home_path}
sshpass ssh ${client3} -i ${client3_cert} ${unzip_raxos}

echo "Client 4"
sshpass ssh ${client4} -i ${client4_cert} ${reset_directory}
sshpass ssh ${client4} -i ${client4_cert} ${kill_insstances}
scp -i ${client4_cert} ${local_zip_path} ${client4}:${replica_home_path}
sshpass ssh ${client4} -i ${client4_cert} ${unzip_raxos}

echo "Client 5"
sshpass ssh ${client5} -i ${client5_cert} ${reset_directory}
sshpass ssh ${client5} -i ${client5_cert} ${kill_insstances}
scp -i ${client5_cert} ${local_zip_path} ${client5}:${replica_home_path}
sshpass ssh ${client5} -i ${client5_cert} ${unzip_raxos}

rm ${local_zip_path}

echo "setup complete"
