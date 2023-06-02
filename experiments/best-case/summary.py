import os
import sys

currentdir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(currentdir)
sys.path.append(parentdir + "/python")
from performance_extract import *

setting = sys.argv[1]  # LAN or WAN
numIter = sys.argv[2]

if setting != "LAN" and setting != "WAN":
    exit("wrong input, input should be LAN/WAN")

if int(numIter) < 4:
    exit("at least 4 iterations needed")

replicaBatchSize = 2000
replicaBatchTime = 4000

if setting == "WAN":
    replicaBatchSize = 3000
    replicaBatchTime = 5000

propTime = 0

if setting == "WAN":
    propTime = 5

iterations = list(range(1, int(numIter) + 1))
arrivals = []

if setting == "LAN":
    arrivals = [1000, 10000, 20000, 30000, 40000, 50000, 80000, 100000, 110000, 112000, 115000, 120000, 130000, 150000,
                180000, 200000]

if setting == "WAN":
    arrivals = [200, 1000, 5000, 10000, 15000, 20000, 25000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 100000]


def getEPaxosSummary():
    l_records = []
    for arrival in arrivals:
        record = ["epaxos-exec", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/epaxos/" + str(arrival) + "/" + str(int(replicaBatchSize)) \
                   + "/" + str(replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/execution/"
            t, l, n, e = getEPaxosPaxosPerformance(root, 7, 5)
            throughput.append(t)
            latency.append(l)
            nine9.append(n)
            err.append(e)
        record.append(int(sum(remove_farthest_from_median(throughput, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(latency, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(nine9, 1)) / (len(iterations) - 1)))
        record.append(int(sum(remove_farthest_from_median(err, 1)) / (len(iterations) - 1)))
        l_records.append(record)

        record = ["epaxos-commit", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/epaxos/" + str(arrival) + "/" + str(int(replicaBatchSize)) \
                   + "/" + str(replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/commit/"
            t, l, n, e = getEPaxosPaxosPerformance(root, 7, 5)
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


def getPaxosSummary():
    l_records = []
    for arrival in arrivals:
        record = ["paxos", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/paxos/" + str(arrival) + "/" + str(int(replicaBatchSize)) \
                   + "/" + str(replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/execution/"
            t, l, n, e = getEPaxosPaxosPerformance(root, 7, 5)
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
            root = "experiments/best-case/logs/quepaxa/" + str(arrival) + "/" + str(int(replicaBatchSize)) + "/" + str(
                replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/" + str(0) + "/" + str(
                propTime) + "/execution/"
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

        record = ["quepaxa-1", str(arrival * 5)]
        throughput, latency, nine9, err = [], [], [], []
        for iteration in iterations:
            root = "experiments/best-case/logs/quepaxa/" + str(arrival) + "/" + str(int(replicaBatchSize)) + "/" + str(
                replicaBatchTime) + "/" + str(setting) + "/" + str(iteration) + "/" + str(1) + "/" + str(
                propTime) + "/execution/"
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

ePaxosSummary = getEPaxosSummary()
paxosSummary = getPaxosSummary()
quePaxaSummary = getQuePaxaSummary()

records = records + ePaxosSummary + paxosSummary + quePaxaSummary

import csv

with open("experiments/best-case/logs/summary.csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)

epaxos_exec_throughput = []
epaxos_exec_latency = []
epaxos_exec_tail = []

epaxos_no_exec_throughput = []
epaxos_no_exec_latency = []
epaxos_no_exec_tail = []

for e in ePaxosSummary:
    if e[0] == "epaxos-commit":
        epaxos_no_exec_throughput.append(e[2])
        epaxos_no_exec_latency.append(e[3])
        epaxos_no_exec_tail.append(e[4])
    elif e[0] == "epaxos-exec":
        if 1000 < float(e[1]) < 6000:
            continue  # arrival rate 5000 requires a low batch size
        epaxos_exec_throughput.append(e[2])
        epaxos_exec_latency.append(e[3])
        epaxos_exec_tail.append(e[4])
    else:
        exit("should not happen")

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
    if ra[0] == "quepaxa-1":
        quepaxa_throughput.append(ra[2])
        quepaxa_latency.append(ra[3])
        quepaxa_tail.append(ra[4])

import matplotlib.pyplot as plt


def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList


# ninty latency

plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()

if setting == "LAN":
    # ax.set_xlim([0, 350])
    ax.set_ylim([0, 80])

if setting == "WAN":
    # ax.set_xlim([0, 390])
    ax.set_ylim([0, 5000])


plt.plot(di_func(quepaxa_throughput), di_func(quepaxa_tail), 'b.-', label="QuePaxa")
plt.plot(di_func(paxos_throughput), di_func(paxos_tail), 'y*-', label="Multi-Paxos")
plt.plot(di_func(epaxos_no_exec_throughput), di_func(epaxos_no_exec_tail), 'cx-', label="Epaxos-commit")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_tail), 'mo-.', label="Epaxos-exec")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('99 percentile Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case/logs/throughput_tail'+setting+'.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()

# median latency

plt.figure(figsize=(5, 4))
plt.rcParams.update({'font.size': 14.30})
ax = plt.gca()
ax.grid()

if setting == "LAN":
    # ax.set_xlim([0, 700])
    ax.set_ylim([0, 7])

if setting == "WAN":
    # ax.set_xlim([0, 360])
    ax.set_ylim([0, 3000])

plt.plot(di_func(quepaxa_throughput), di_func(quepaxa_latency), 'b.-', label="QuePaxa")
plt.plot(di_func(paxos_throughput), di_func(paxos_latency), 'y*-', label="Multi-Paxos")
plt.plot(di_func(epaxos_no_exec_throughput), di_func(epaxos_no_exec_latency), 'cx-', label="Epaxos-commit")
plt.plot(di_func(epaxos_exec_throughput), di_func(epaxos_exec_latency), 'mo-', label="Epaxos-exec")

plt.xlabel('Throughput (x 1k cmd/sec)')
plt.ylabel('Median Latency (ms)')
plt.legend()
plt.savefig('experiments/best-case/logs/throughput_median'+setting+'.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()