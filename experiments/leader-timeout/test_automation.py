import os
import sys

numIter = sys.argv[1]

if int(numIter) < 4:
    exit("at least 4 iterations needed")

os.system("/bin/bash experiments/setup-5/setup.sh")

leaderTimeouts = [1000, 20000, 50000, 100000, 150000, 200000, 250000, 300000, 350000, 400000, 450000, 500000, 550000,
                  600000, 650000, 750000]

iterations = list(range(1, int(numIter) + 1))
modes = ["performance", "recovery"]

for iteration in iterations:
    for leaderTimeout in leaderTimeouts:
        for mode in modes:
            os.system(
                "/bin/bash experiments/leader-timeout/quepaxa.sh " + str(iteration) + " "
                + str(leaderTimeout) + " "
                + str(mode) )

            for algo in ["paxos", "raft"]:
                os.system(
                    "/bin/bash experiments/leader-timeout/paxos_raft.sh " + str(algo) + " "
                    + str(iteration) + " "
                    + str(leaderTimeout) + " "
                    + str(mode) )
