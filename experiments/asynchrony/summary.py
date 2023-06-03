import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir + "/python")
from performance_extract import *

numIter = sys.argv[1]
device = sys.argv[2]

if int(numIter) < 4:
    exit("at least 4 iterations needed")

iterations = list(range(1, int(numIter) + 1))
arrivals = [200, 1000, 2000, 3000, 5000, 6000, 7000, 8000, 9000, 10000, 12000]

def getPaxosSummary():
    l_records = []
    for arrival in arrivals:
        record = ["paxos", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/asynchrony/logs/paxos/" + str(arrival) + "/" + str(int(iteration)) \
                   + "/" + str(device) + "/"
            t, l, n, e = getPaxosRaftPerformance(root, 21, 5)
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


def getRaftSummary():
    l_records = []
    for arrival in arrivals:
        record = ["raft", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/asynchrony/logs/raft/" + str(arrival) + "/" + str(int(iteration)) \
                   + "/" + str(device) + "/"
            t, l, n, e = getPaxosRaftPerformance(root, 21, 5)
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


def getQuePaxaSummary():
    l_records = []
    for arrival in arrivals:
        record = ["quepaxa-0", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/asynchrony/logs/quepaxa/" + str(arrival) + "/" + str(int(iteration)) + "/" + str(
                device) + "/"
            t, l, n, e = getQuePaxaPerformance(root, 21, 5)
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


headers = ["algo", "arrivalRate", "throughput", "median latency", "99%", "error rate"]
records = [headers]

raftSummary = getRaftSummary()
paxosSummary = getPaxosSummary()
quePaxaSummary = getQuePaxaSummary()

records = records + raftSummary + paxosSummary + quePaxaSummary

import csv

with open("experiments/asynchrony/logs/summary.csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)

raft_throughput = []
raft_latency = []
raft_tail = []

for r in raftSummary:
    raft_throughput.append(r[2])
    raft_latency.append(r[3])
    raft_tail.append(r[4])

paxos_throughput = []
paxos_latency = []
paxos_tail = []

for p in paxosSummary:
    paxos_throughput.append(p[2])
    paxos_latency.append(p[3])
    paxos_tail.append(p[4])

quepaxa_throughput = []
quepaxa_latency = []
quepaxa_tail = []

for ra in quePaxaSummary:
    quepaxa_throughput.append(ra[2])
    quepaxa_latency.append(ra[3])
    quepaxa_tail.append(ra[4])

import matplotlib.pyplot as plt


def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList

# median latency


plt.figure(figsize=(5, 1.8))
plt.rcParams.update({'font.size': 10.00})
ax = plt.gca()
ax.grid()
# ax.set_xlim([0, 350])
ax.set_ylim([0, 1100])

plt.plot(di_func(quepaxa_throughput), di_func(quepaxa_latency), 'b.-', label="QuePaxa")
plt.plot(di_func(paxos_throughput), di_func(paxos_latency), 'y*-', label="Multi-Paxos")
plt.plot(di_func(raft_throughput), di_func(raft_latency), 'gx-', label="Raft")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('Median latency (ms)')
plt.legend()
plt.savefig('experiments/asynchrony/logs/asynchrony.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()