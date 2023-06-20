import sys

numIter = sys.argv[1]

if int(numIter) < 4:
    exit("at least 4 iterations needed")

iterations = list(range(1, int(numIter) + 1))
leaderTimeouts = [1000, 20000, 50000, 100000, 150000, 200000, 250000, 300000, 350000, 400000, 450000, 500000, 550000,
                  600000, 650000, 750000]


def steps_per_slot(file_name):
    try:
        f = open(file_name, 'r')
    except OSError:
        exit("Error in " + file_name + "\n")
    with f:
        content = f.readlines()
    if len(content) < 2:
        exit("Error in " + file_name + "\n")
    for c in content:
        if c.startswith("Average number of steps"):
            return float(c.strip().split(" ")[6])
    exit("should not happen")


def get_que_paxa_summary():
    l_records = []
    for timeout in leaderTimeouts:
        record = [str(timeout)]
        steps = 0
        for iteration in iterations:
            file_name = "experiments/leader-timeout/logs/performance/quepaxa/" + str(iteration) + "/" + str(
                int(timeout)) + "/1.log"
            steps = steps + steps_per_slot(file_name)

        record.append(steps / (len(iterations)))
        l_records.append(record)
    return l_records


quePaxaSummary = get_que_paxa_summary()
for s in quePaxaSummary:
    print("timeout " + str(s[0]) + " average number of steps per slot " + str(s[1]))
