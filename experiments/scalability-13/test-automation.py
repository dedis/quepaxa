import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir + "/python")
from performance_extract import *

numIter = sys.argv[1]

if int(numIter) < 3:
    exit("at least 3 iterations needed")

os.system("/bin/bash experiments/setup-13/setup.sh")

iterations = list(range(1, int(numIter) + 1))
arrivals = [467337, 380000]

for iteration in iterations:
    for arrival in arrivals:
        os.system(
            "/bin/bash experiments/scalability-13/quepaxa.sh " + str(int(arrival)) + " "
            + str(iteration))


def getQuePaxaSummary():
    l_records = []
    for arrival in arrivals:
        record = ["quepaxa", str(arrival * 13)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/scalability/logs/quepaxa/" + str(arrival) + "/" + str(int(iteration)) + "/"
            t, l, n, e = getQuePaxaPerformance(root, 21, 13)
            throughput.append(t)
            latency.append(l)
            nine9.append(n)
            err.append(e)
        record.append(int(sum(remove_farthest_from_median(throughput, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(latency, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(nine9, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(err, 1)) / (len(iterations) - 1)))
        l_records.append(record)
    return l_records

print(getQuePaxaSummary())