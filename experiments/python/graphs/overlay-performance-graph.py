import matplotlib.pyplot as plt

throughput_overlay = [5001, 25014, 49984, 75008, 125094, 274993, 374954, 500040, 550055, 599936]
median_latency_overlay = [554, 792.8, 960.2, 1260, 1878.4, 2923, 4661.8, 10840, 14872.8, 18774.4]
ninetynine_latency_overlay = [6152.6, 49942.6, 49985.4, 75357.4, 89009.8, 130943, 287972.7, 651230.7, 742857.1,
                              1176964.1]

throughput_basic_raxos = [4990, 25017, 49997, 99973, 200012, 299882, 399850, 499909, 549915]
median_latency_raxos = [1938.2, 5224.8, 11551.4, 29241.5, 42868, 100261.4, 210789.6, 800948.2, 1249306.2]
ninetynine_latency_raxos = [73599, 561568.5, 663196.9, 665300.9, 931013.6, 1302786, 3496800.7, 5591776.3, 5991776.3]
# plt.rcParams.update({'font.size': 13.45})

ax = plt.gca()
# ax.set_xscale('log')
# ax.set_xlim([0, 749865 + 10])
# ax.set_ylim([0, 14200])

# plt.plot(throughput_overlay, median_latency_overlay, 'r*-', label="Overlay")
# plt.plot(throughput_basic_raxos, median_latency_raxos, 'g*-', label="QSCOD")
plt.plot(throughput_overlay, ninetynine_latency_overlay, 'r*-', label="Overlay")
plt.plot(throughput_basic_raxos, ninetynine_latency_raxos, 'g*-', label="QSCOD")

plt.xlabel('Throughput (requests per second)')
# plt.ylabel('Median latency (micro seconds)')
plt.ylabel('Tail latency 99% (micro seconds)')
plt.legend()

# plt.savefig('/home/pasindu/Desktop/Test/median_latency.png')
plt.savefig('/home/pasindu/Desktop/Test/tail_latency.png')
