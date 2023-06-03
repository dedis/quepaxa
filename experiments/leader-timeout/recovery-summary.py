import os
import sys
from threading import Barrier, Thread

numIter = sys.argv[1]

if int(numIter) < 4:
    exit("at least 4 iterations needed")

leaderTimeouts = [1000, 20000, 50000, 100000, 150000, 200000, 250000, 300000, 350000, 400000, 450000, 500000, 550000,
                  600000, 650000, 750000, 850000, 950000, 1050000, 1150000, 1250000, 1350000, 1450000]
leaderTimeouts.reverse()


def execute(timeout_l, barrier):
    os.system("python3 experiments/leader-timeout/timeout_summary.py " + numIter + " " + str(timeout_l))
    barrier.wait()
    print("finished " + str(timeout_l))


sys.stdout.flush()
barrier = Barrier(len(leaderTimeouts) + 1)

for timeout in leaderTimeouts:
    worker = Thread(target=execute, args=(timeout, barrier))
    worker.start()

barrier.wait()
print("finished!")

string = "cat "
for leaderTimeout in leaderTimeouts:
    string = string + " experiments/leader-timeout/logs/longest_zero_local" + str(leaderTimeout) + ".csv"

os.system(string + " > " + "experiments/leader-timeout/logs/recovery_time_summary.csv")


import csv

with open('experiments/leader-timeout/logs/recovery_time_summary.csv', newline='') as csvfile:
    data = list(csv.reader(csvfile))

quepaxa_timeouts = []
quepaxa_recovery = []

paxos_timeouts = []
paxos_recovery = []

raft_timeouts = []
raft_recovery = []

for row in data:
    if row[0] == "quepaxa":
        quepaxa_timeouts.append(int(row[1]))
        quepaxa_recovery.append(int(row[2]))
    if row[0] == "paxos":
        paxos_timeouts.append(int(row[1]))
        paxos_recovery.append(int(row[2]))
    if row[0] == "raft":
        raft_timeouts.append(int(row[1]))
        raft_recovery.append(int(row[2]))

# graph

import matplotlib.pyplot as plt
import matplotlib


def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList


plt.rcParams.update({'font.size': 9.00})
fig1, ax1 = plt.subplots(figsize=(5, 1.5))

ax1.plot(di_func(quepaxa_timeouts), quepaxa_recovery, 'b.-', label="QuePaxa")
ax1.plot(di_func(paxos_timeouts), paxos_recovery, 'y*-', label="Multi-Paxos")
ax1.plot(di_func(raft_timeouts), raft_recovery, 'gx-', label="Raft")
ax1.axvline(x = 179.83, linestyle='dotted',  color = 'm', label = 'Round trip latency')

ax1.set_xscale('log')
ax1.set_yscale('log')
ax1.grid()
ax1.set_xticks([100, 200, 300, 500])
ax1.get_xaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())
# ax1.set_yticks([100, 1000,25000])
ax1.get_yaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())


plt.xlabel('Leader Timeout / Hedging Delay (ms)')
plt.ylabel('Recovery time \n(ms)')
plt.legend()
plt.savefig('experiments/leader-timeout/logs/timeout_recovery.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()

