import sys

numIter = sys.argv[1]
timeout = sys.argv[2]

if int(numIter) < 3:
    exit("at least 3 iterations needed")

iterations = list(range(1, int(numIter) + 1))


def findLongestZero(x):
    longest = 0
    start = 0
    start_i = 0
    end_i = len(x) - 1
    while True:
        count = 0
        start_i_l = start
        end_i_l = start
        while start < len(x) and x[start] == 0:
            start = start + 1
            count = count + 1
            end_i_l = start
        if count > longest:
            longest = count
            start_i = start_i_l
            end_i = end_i_l
        while start < len(x) and x[start] != 0:
            start = start + 1
        if start >= len(x):
            return longest, start_i, end_i - 1


def binToMiiliSeconds(times):  # times is a microseconds array
    times = list(dict.fromkeys(times))
    times.sort()
    x = list(range(1, 60000, 1))
    time_x = []
    prev = 0
    index = 0
    for i in x:
        count = 0
        index_l = index
        for t in times[index_l:]:
            if prev * 1000 < t < i * 1000:
                count = count + 1
            if t >= i * 1000:
                break
            index = index + 1
            if index >= (len(times) - 1):
                break
        time_x.append(count)
        prev = i

    return time_x


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
            times.append(end)

    if len(times) == 0:
        return []

    return list(set(times))


def getTimes(algo, start_client, num_client):
    times = []
    for iteration in iterations:
        local_root = "experiments/leader-timeout/logs/recovery/"+ algo + "/" + str(iteration) + "/"+str(timeout)+"/"
        for cl in list(range(start_client, start_client + num_client, 1)):
            file_name = local_root + str(cl) + ".txt"
            times_local = extractTimes(file_name)
            if len(times_local) > 0:
                for t in times_local:
                    times.append(t)
    return times


def execute():

    records = []

    # paxos
    times = getTimes("paxos", 21, 5)
    if len(times) > 0:
        bins = binToMiiliSeconds(times)
        longestZero = findLongestZero(bins)
        records.append(["paxos", str(timeout), longestZero[0], longestZero[1], longestZero[2]])

    # raft
    times = getTimes("raft", 21, 5)
    if len(times) > 0:
        bins = binToMiiliSeconds(times)
        longestZero = findLongestZero(bins)
        records.append(["raft", str(timeout), longestZero[0], longestZero[1], longestZero[2]])

    # quepaxa
    times = getTimes("quepaxa", 21, 5)
    if len(times) > 0:
        bins = binToMiiliSeconds(times)
        longestZero = findLongestZero(bins)
        records.append(["quepaxa", str(timeout), longestZero[0], longestZero[1], longestZero[2]])
    return records


records_timeout = execute()

import csv
with open("experiments/leader-timeout/logs/longest_zero_local" + sys.argv[2] + ".csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records_timeout)
