import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir + "/python")
from performance_extract import *

numIter = sys.argv[1]

if int(numIter) < 4:
    exit("at least 4 iterations needed")

iterations = list(range(1, int(numIter) + 1))
leaderTimeouts = [1000, 20000, 50000, 100000, 150000, 200000, 250000, 300000, 350000, 400000, 450000, 500000, 550000,
                  600000, 650000, 750000, 850000, 950000, 1050000, 1150000, 1250000, 1350000, 1450000]

def getPaxosSummary():
    l_records = []
    for timeout in leaderTimeouts:
        record = ["paxos", str(timeout)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/leader-timeout/logs/performance/paxos/" + str(iteration) + "/" + str(int(timeout))
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
    for timeout in leaderTimeouts:
        record = ["raft", str(timeout)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/leader-timeout/logs/performance/raft/" + str(iteration) + "/" + str(int(timeout))+ "/"
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
    for timeout in leaderTimeouts:
        record = ["quepaxa", str(timeout)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/leader-timeout/logs/performance/quepaxa/" + str(iteration) + "/" + str(int(timeout)) + "/"
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


headers = ["algo", "timeout", "throughput", "median latency", "99%", "error rate"]
records = [headers]

raftSummary = getRaftSummary()
paxosSummary = getPaxosSummary()
quePaxaSummary = getQuePaxaSummary()

records = records + raftSummary + paxosSummary + quePaxaSummary

import csv

with open("experiments/leader-timeout/logs/summary.csv", "w", newline="") as f:
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


plt.rcParams.update({'font.size': 9.00})
fig1, ax1 = plt.subplots(figsize=(5, 1.5))


ax1.plot(di_func(leaderTimeouts), quepaxa_throughput, 'b.-', label="QuePaxa")
ax1.plot(di_func(leaderTimeouts), paxos_throughput, 'y*-', label="Multi-Paxos")
ax1.plot(di_func(leaderTimeouts), raft_throughput, 'gx-', label="Raft")
ax1.axvline(x = 179.83, linestyle='dotted',  color = 'm', label = 'Round trip latency')

ax1.set_xscale('log')
ax1.set_yscale('log')
ax1.grid()
ax1.set_xticks([100, 200, 300, 500])
ax1.get_xaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())
ax1.set_yticks([1000,4000, 10000,25000])
ax1.get_yaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())

plt.xlabel('Leader Timeout / Hedging Delay (ms)')
plt.ylabel('Throughput \n cmd/sec')
plt.legend()
plt.savefig('aws/timeout-performance/timeout_performance.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()