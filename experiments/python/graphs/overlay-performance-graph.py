import matplotlib.pyplot as plt

throughput = [5001,25014, 49984,75008, 125094,274993,374954,500040,550055,599936]
median_latency = [554,792.8,960.2,1260,1878.4,2923,4661.8,10840,14872.8,18774.4]
ninetynine_latency = [6152.6,49942.6,49985.4,75357.4,89009.8,130943,287972.7,651230.7,742857.1,1176964.1]

# plt.rcParams.update({'font.size': 13.45})

ax = plt.gca()
# ax.set_xscale('log')
# ax.set_xlim([0, 749865 + 10])
# ax.set_ylim([0, 14200])

plt.plot(throughput, ninetynine_latency, 'r*-', label="QSCOD")

plt.xlabel('Throughput (requests per second)')
plt.ylabel('99% latency (micro seconds)')
plt.legend()

plt.savefig('/home/oem/Desktop/Test/tail_latency..png')
