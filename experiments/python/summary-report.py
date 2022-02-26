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

for arrivalrate in [1000, 5000, 10000, 15000, 25000, 55000, 75000, 100000, 110000, 120000]:
    newRecord = [str(arrivalrate)]
    throughputs = []
    latency = []
    tail_latency = []
    errors = []

    for client in range(5, num_clients + 5):
        file_name = "/home/oem/Desktop/Test/" + str(arrivalrate) + "/logs/" + str(client) + ".log"
        with open(file_name) as f:
            content = f.readlines()
        if not (content[10].strip().split(" ")[0] == "Throughput" and content[11].strip().split(" ")[0] == "Median"):
            print("Error in " + file_name + "\n")
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

with open("/home/oem/Desktop/Test/best_case_tail.csv", "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerows(records)
