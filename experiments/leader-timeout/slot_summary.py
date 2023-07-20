import sys

numIter = sys.argv[1]

iterations = list(range(1, int(numIter) + 1))
leaderTimeouts = [50000, 100000, 150000, 200000, 250000, 300000, 350000, 400000, 450000, 500000, 550000,
                  600000, 650000, 750000]


def slots_and_steps(file_name):
    try:
        f = open(file_name, 'r')
    except OSError:
        exit("Error in " + file_name + "\n")
    with f:
        content = f.readlines()

    for c in content:
        if c.startswith("Average number of steps per slot"):
            return [float(c.strip().split(" ")[9][:-1]), float(c.strip().split(" ")[12])]
    exit(file_name + " error")


def get_que_paxa_summary():
    print("timeout", "slots", "steps", "average")
    l_records = []
    for timeout in leaderTimeouts:
        record = [str(timeout)]
        slots = 0
        steps = 0
        for iteration in iterations:
            for cl in [1, 2, 3, 4, 5]:
                file_name = "experiments/leader-timeout/logs/performance/quepaxa/" + str(iteration) + "/" + str(
                    int(timeout)) + "/" + str(cl) + ".log"
                file_stat = slots_and_steps(file_name)
                if len(file_stat) == 2:
                    slots = slots + file_stat[0]
                    steps = steps + file_stat[1]
                else:
                    exit("error in file name "+file_name)

        record.append(steps / slots)
        l_records.append(record)
        print(str(timeout) + " " + str(slots) + " " + str(steps) + " " + str(steps / slots))
    return l_records


quePaxaSummary = get_que_paxa_summary()

quePaxaAverage = []

for s in quePaxaSummary:
    quePaxaAverage.append(s[1])

import matplotlib.pyplot as plt
import matplotlib


def di_func(array):
    returnList = []
    for l in array:
        returnList.append(l / 1000)
    return returnList


plt.rcParams.update({'font.size': 9.00})
fig1, ax1 = plt.subplots(figsize=(5, 1.5))

ax1.plot(di_func(leaderTimeouts), quePaxaAverage, 'g.-', label="QuePaxa")
ax1.axvline(x=179.83, linestyle='dotted', color='m', label='Round trip latency')

ax1.set_xscale('log')
ax1.grid()
ax1.set_xticks([50, 100, 300, 500])
ax1.get_xaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())
ax1.set_yticks([1,2, 3,4,5])
ax1.get_yaxis().set_major_formatter(matplotlib.ticker.ScalarFormatter())

plt.xlabel('Leader Timeout / Hedging Delay (ms)')
plt.ylabel('Steps per slot')
plt.legend()
plt.savefig('experiments/leader-timeout/logs/steps.pdf', bbox_inches='tight', pad_inches=0)
plt.close()
plt.clf()
plt.cla()
