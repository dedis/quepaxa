import os
import sys

numIter = sys.argv[1]

if int(numIter) < 4:
    exit("at least 4 iterations needed")

os.system("/bin/bash experiments/setup-5/setup.sh")

iterations = list(range(1, int(numIter) + 1))
arrivals = [200, 1000, 2000, 3000, 5000, 6000, 7000, 8000, 9000, 10000, 12000]
arrivals.reverse()

for iteration in iterations:
    for arrival in arrivals:
        os.system(
            "/bin/bash experiments/asynchrony/quepaxa.sh " + str(arrival) + " "
            + str(iteration))

        for algo in ["paxos", "raft"]:
            os.system(
                "/bin/bash experiments/asynchrony/paxos_raft.sh " + str(algo) + " "
                + str(arrival) + " "
                + str(iteration))
