import sys

num_clients = int(sys.argv[1])

headers = ["Arrival Rate", "Error Rate", "Median Latency (micro seconds)", "99% latency (micro seconds)",
           "Throughput (requests per second)"]

for i in range(1, num_clients + 1):
    headers.append("Client " + str(i) + " Throughput (requests per second)")

for i in range(1, num_clients + 1):
    headers.append("Client " + str(i) + " Median Latency (micro seconds)")

# best case
records = [headers]

for arrivalrate in [100, 200, 500, 600, 700, 800, 1000, 1200, 1500, 1800, 2000, 2200, 2500, 3000, 3500, 4000, 4500, 5000, 5500, 6000, 6500, 7000]:
    newRecord = [str(arrivalrate)]
    throughputs = []
    latency = []
    tail_latency = []
    errors = []

    for client in range(5, num_clients + 5):
        file_name = "/home/pasindu/Desktop/Test/Non_Batch/" + str(arrivalrate) + "/logs/" + str(client) + ".log"
        with open(file_name) as f:
            content = f.readlines()

        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]

        if not (content[10].strip().split(" ")[0] == "Throughput" and content[11].strip().split(" ")[0] == "Median"):
            print("Error in " + file_name + "\n")
            continue
        throughputs.append(float(content[10].strip().split(" ")[5]))
        latency.append(float(content[11].strip().split(" ")[6]))
        tail_latency.append(float(content[12].strip().split(" ")[7]))
        errors.append(float(content[13].strip().split(" ")[3]))

    avgLatency = sum(latency) / num_clients
    tailLatency = sum(tail_latency) / num_clients
    avgThroughput = sum(throughputs)
    errorRate = sum(errors)

    newRecord.append(str(errorRate))
    newRecord.append(str(avgLatency))
    newRecord.append(str(tailLatency))
    newRecord.append(str(avgThroughput))
    for t in throughputs:
        newRecord.append(t)
    for l in latency:
        newRecord.append(l)

    records.append(newRecord)

import csv

with open("/home/pasindu/Desktop/Test/Non_Batch/non_leader.csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)
