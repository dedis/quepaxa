import matplotlib.pyplot as plt
import matplotlib

def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList

paxos_recovery = [ 294, 199, 170, 164, 17321, 19321, 21194]
paxos_timeouts = [600000, 550000, 500000, 450000, 400000, 350000, 300000]

raft_recovery = [166, 181, 195, 186, 163, 193, 156, 27917, 26682 ]
raft_timeouts = [600000, 550000, 500000, 450000, 400000, 350000, 300000, 250000, 200000]

raxos_recovery = [143, 179, 243, 147, 124, 141, 225, 156, 146, 212, 215, 338]
raxos_timeouts = [600000, 550000, 500000, 450000, 400000, 350000, 300000, 250000, 200000, 150000, 100000, 50000]

timeouts = [50000, 100000, 150000, 200000, 250000, 300000, 350000, 400000, 450000, 500000, 550000, 600000]
paxos_throughput = [0, 5.5, 700, 3128.5, 10007, 24735, 24746, 25010, 25013, 25009, 25004, 25006]
raft_throughput = [0, 0, 1486, 3192, 24971, 25005, 24969, 24979, 24990, 24989, 25002, 25017]
raxos_throughput = [25105, 25091, 25096, 25099, 25102, 25099, 25097, 25094, 25098, 25100, 25095, 25097]

raxos_slots = [2.822, 1.0674, 1,1,1,1,1,1,1,1,1,1]

plt.rcParams.update({'font.size': 11.00})
fig, (ax1, ax2, ax3) = plt.subplots(3, 1)
plt.subplots_adjust(wspace=0.03)
fig.set_figheight(5)
fig.set_figwidth(5)

ax1.plot(di_func(timeouts), raxos_slots, 'b.-', label="QuePaxa")
ax1.axvline(x = 179.83, linestyle='dotted',  color = 'm')

ax1.set_xscale('log')
ax1.set_xticks([50, 100,200, 300, 500])
ax1.set_xticklabels([])
ax1.set_yticks([1,2,3])
ax1.get_yaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())
ax1.grid()
ax1.set_ylabel("Average number \n of steps per \n slot")


ax2.plot(di_func(raxos_timeouts), raxos_recovery, 'b.-')
ax2.plot(di_func(paxos_timeouts), paxos_recovery, 'y*-')
ax2.plot(di_func(raft_timeouts), raft_recovery, 'gx-')
ax2.axvline(x = 179.83, linestyle='dotted',  color = 'm')

ax2.set_xscale('log')
ax2.set_yscale('log')
ax2.set_xticks([50, 100, 200, 300, 500])
ax2.set_xticklabels([])
ax2.grid()

ax2.set_yticks([100, 300, 1000, 3000])
ax2.set_ylim([100, 3000])
ax2.get_yaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())
ax2.set_ylabel("Recovery time \n (ms)")

ax3.plot(di_func(timeouts), raxos_throughput, 'b.-', label="QuePaxa")
ax3.plot(di_func(timeouts), paxos_throughput, 'y*-', label="Multi-Paxos")
ax3.plot(di_func(timeouts), raft_throughput, 'gx-', label="Raft")
ax3.axvline(x = 179.83, linestyle='dotted',  color = 'm', label = 'Round trip latency')

ax3.set_xscale('log')
ax3.set_yscale('log')
ax3.set_xticks([50, 100,200, 300, 500])
ax3.get_xaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())
ax3.grid()
ax3.set_yticks([5, 1000,4000, 25000])
ax3.get_yaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())
ax3.set_ylabel("Throughput \n (cmd/sec)")


plt.xlabel('Leader Timeout / Hedging Delay (ms)')
plt.legend(prop={'size': 8},loc="lower right")
plt.savefig('experiments/leader-timeout/logs/hedging.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()