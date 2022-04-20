import matplotlib.pyplot as plt

throughput_overlay = [5001, 25014, 49984, 75008, 125094, 274993, 374954, 500040, 550055, 599936]
median_latency_overlay = [554, 792.8, 960.2, 1260, 1878.4, 2923, 4661.8, 10840, 14872.8, 18774.4]
ninetynine_latency_overlay = [6152.6, 49942.6, 49985.4, 75357.4, 89009.8, 130943, 287972.7, 651230.7, 742857.1,
                              1176964.1]

throughput_leader_raxos = [4990, 25017, 49997, 99973, 200012, 299882, 399850, 499909, 549915]
median_latency_leader_raxos = [1938.2, 5224.8, 11551.4, 29241.5, 42868, 100261.4, 210789.6, 800948.2, 1249306.2]
ninetynine_latency_leader_raxos = [73599, 561568.5, 663196.9, 665300.9, 931013.6, 1302786, 3496800.7, 5591776.3, 5991776.3]

throughput_non_leader_raxos = [4987, 24969, 49973, 100050, 200011, 299930]
median_latency_non_leader_raxos = [47219.4, 20166556.1, 25202255.9, 30129535, 47906241.8, 59937966.8]
ninetynine_latency_non_leader_raxos = [511607.9, 36956067.5, 40790758, 48333083.9, 55868651.2,61393796.9]

throughput_non_leader_no_batch_raxos=[497, 1001, 2499, 2992, 3502, 3998]
median_latency_non_leader_no_batch_raxos=[3213.5, 4027.1, 10124159, 36629438.9, 59014385.9, 73076252.7]
ninetynine_latency_non_leader_no_batch_raxos=[26266.1, 69434.7, 24069363, 60387052.7, 82837442.1, 97468506]


plt.rcParams.update({'font.size': 13.45})

ax = plt.gca()
ax.set_xscale('log')
# ax.set_xlim([0, 749865 + 10])
# ax.set_ylim([0, 14200])

# plt.plot(throughput_overlay, median_latency_overlay, 'r*-', label="Overlay")
# plt.plot(throughput_leader_raxos, median_latency_leader_raxos, 'g*-', label="QSCOD-Leader")
# plt.plot(throughput_non_leader_raxos, median_latency_non_leader_raxos, 'b*-', label="QSCOD-Non-Leader")

# plt.plot(throughput_overlay, ninetynine_latency_overlay, 'r*-', label="Overlay")
# plt.plot(throughput_leader_raxos, ninetynine_latency_leader_raxos, 'g*-', label="QSCOD-Leader")
# plt.plot(throughput_non_leader_raxos, ninetynine_latency_non_leader_raxos, 'b*-', label="QSCOD-Non-Leader")

plt.plot(throughput_non_leader_raxos, median_latency_non_leader_raxos, 'r-', label="Batch size 50")
plt.plot(throughput_non_leader_no_batch_raxos, median_latency_non_leader_no_batch_raxos, 'b-', label="Batch size 1")


plt.xlabel('Throughput (requests per second)')
plt.ylabel('Median latency (micro seconds)')
# plt.ylabel('Tail latency 99% (micro seconds)')
plt.legend()

# plt.savefig('/home/pasindu/Desktop/Test/median_latency.png')
# plt.savefig('/home/pasindu/Desktop/Test/tail_latency.png')
plt.savefig('/home/pasindu/Desktop/Test/batching.png')
