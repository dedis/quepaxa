ip=$1
dport=$2
sshpass ssh ubuntu@${ip} "sudo tc qdisc add dev ens5 root handle 1: prio priomap 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 ; sudo tc qdisc add dev ens5 parent 1:2 handle 20: netem delay 4000ms ; sudo tc filter add dev ens5 parent 1:0 protocol ip u32 match ip dport ${dport} 0xffff flowid 1:2 ; sleep 3 ; sudo tc qdisc del dev ens5 root"