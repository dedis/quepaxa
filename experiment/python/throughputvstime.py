import matplotlib.pyplot as plt
import sys

def getLatency(file_path, metric_local):
    content = []
    max_iterations = 1
    for iteration in range(1, max_iterations+1, 1):
        for clt in range(21, 22, 1):
            file_name = file_path
            with open(file_name) as f:
                lines = f.readlines()
                for line in lines:
                    content.append(line)

    x = range(1, 61, 1)
    time_x = []
    y = []
    times = []
    for c in content:
        item = c.strip().split(",")
        times.append([int(item[1])/1000, int(item[2])/1000])

    prev = 0
    for i in x:
        latencySum = 0
        latencyCount = 0
        for t in times:
            if t[1] > prev*1000 and t[1]< i*1000:
                latencyCount = latencyCount+1
                latencySum = latencySum + t[1] - t[0]


        if latencyCount > 0:
            time_x.append(i)
            if metric_local == "latency":
                y.append(latencySum/latencyCount)
            if metric_local == "throughput":
                y.append(latencyCount)
        prev = i

    return time_x,y

metric = sys.argv[1]

x, y = getLatency("/home/pasindu/Documents/Raxos/logs/21.txt", metric)

plt.rcParams.update({'font.size': 14.45})
plt.rcParams["figure.figsize"] = (8,6)

ax = plt.gca()
ax.set_xlim([0, 65])

plt.plot(x, y, 'r*-', label = "Raxos")

plt.xlabel('Time (s)')
if metric == "latency":
    plt.ylabel('Median Latency (ms)')
if metric == "throughput":
    plt.ylabel('Throughput (requests per second)')
plt.legend()
plt.grid()

plt.savefig("logs/_"+metric+".png")