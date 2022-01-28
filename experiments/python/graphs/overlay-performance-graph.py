import matplotlib.pyplot as plt

throughput = [4990, 24951, 50012, 75062, 125037, 274882, 374842, 499818, 549932, 600081, 649972, 700021, 749865]
median_latency = [548.8, 773.2, 919.8, 1171, 1720.4, 2460.9, 3235.6, 3813.6, 4224.8, 5308.8, 6391, 7569, 14091.4]
ninetynine_latency = [6497.5, 42675.2, 36287.3, 61275.3, 63253.3, 113650, 192097.4, 222879.5, 137819, 457330.6,
                      387007.8, 510549.8, 1315137.5]

# plt.rcParams.update({'font.size': 13.45})

ax = plt.gca()
# ax.set_xscale('log')
# ax.set_xlim([0, 749865 + 10])
# ax.set_ylim([0, 14200])

plt.plot(throughput, median_latency, 'r*-', label="QSCOD")

plt.xlabel('Throughput (requests per second)')
plt.ylabel('Median Latency (micro seconds)')
plt.legend()

plt.savefig('/home/oem/Desktop/Test/median_latency..png')
