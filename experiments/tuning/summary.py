import matplotlib.pyplot as plt
import sys

numIter = sys.argv[1]
arrival = sys.argv[2] # invoke separately for each arrival rate

if int(numIter) < 4:
    exit("at least 4 iterations needed")

iterations = list(range(1, int(numIter) + 1))


def extractTimes(file_name):
    content = []
    with open(file_name) as f:
        lines = f.readlines()
        for line in lines:
            content.append(line)
    if len(content) == 0:
        return []

    times = []
    for c in content:
        item = c.strip().split(",")
        start = int(item[1])
        end = int(item[2])
        if end - start != 2000000:  # not a timed out request
            times.append([start, end])

    return times


def getTimes(file_root, start_client, num_client):
    times = []
    for iteration in iterations:
        local_root = file_root + str(iteration) + "/"+str(arrival)+"/"
        for cl in list(range(start_client, start_client + num_client, 1)):
            file_name = local_root + str(cl) + ".txt"
            times_local = extractTimes(file_name)
            if len(times_local) > 0:
                for t in times_local:
                    times.append(t)
    return times


def getStats(times):
    times_s = list(range(1, 60000, 10))
    throughputs = []
    latency = []
    prev = 0
    for i in times_s:
        count = 0
        summa = 0
        for t in times:
            if prev * 1000 < t[1] < i * 1000:
                count = count + 1
                summa = summa + (t[1] - t[0])

        throughputs.append(count / len(iterations))
        latency.append(int((summa) / count))
        prev = i

    return [times_s, throughputs, latency]


paxos_s, paxos_t, paxos_l = getStats(getTimes("experiments/tuning/logs/paxos/", 25, 1))

raft_s, raft_t, raft_l = getStats(getTimes("experiments/tuning/logs/raft/", 25, 1))

quepaxa_s, quepaxa_t, quepaxa_l = getStats(getTimes("experiments/tuning/logs/quepaxa/", 25, 1))


def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList


# latency

plt.figure(figsize=(5, 1.5))
plt.rcParams.update({'font.size': 10.00})
ax = plt.gca()
ax.grid()
ax.set_xlim([0, 25])
ax.set_ylim([2, 9])

plt.plot(quepaxa_s, di_func(quepaxa_l), 'b.-', label="QuePaxa")
plt.plot(paxos_s, di_func(paxos_l), 'y*-', label="Multi-Paxos")
plt.plot(raft_s, di_func(raft_l), 'gx-', label="Raft")

plt.xlabel('Time (s)')
plt.ylabel('Average Latency \n (ms)')
plt.legend(ncols=3)
plt.savefig('experiments/tuning/time_latency.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()
