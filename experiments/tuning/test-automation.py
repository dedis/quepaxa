import os
import sys

numIter = sys.argv[1]

if int(numIter) < 4:
    exit("at least 4 iterations needed")

os.system("/bin/bash experiments/setup-5/setup.sh")

iterations = list(range(1, int(numIter) + 1))
arrivals = [100000]

for iteration in iterations:
    for arrival in arrivals:
        for algo in ["paxos", "raft"]:
            os.system(
                "/bin/bash experiments/tuning/paxos_raft.sh " + str(algo) + " "
                + str(iteration) + " "
                + str(arrival))

        os.system(
            "/bin/bash experiments/tuning/quepaxa.sh " + str(iteration) + " "
            + str(arrival))