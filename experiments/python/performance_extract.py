def getPaxosRaftPerformance(root, initClient, numClients):
    throughputs = []
    medians = []
    ninety9s = []
    errors = []
    for cl in list(range(initClient, initClient+numClients, 1)):
        file_name = root + str(cl) + ".log"
        # print(file_name + "\n")
        try:
            f = open(file_name, 'r')
        except OSError:
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        with f:
            content = f.readlines()
        if len(content) < 8:
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]
        if not (content[4].strip().split(" ")[0] == "Throughput" and content[5].strip().split(" ")[0] == "Median"):
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        throughputs.append(float(content[4].strip().split(" ")[5]))
        medians.append(float(content[5].strip().split(" ")[3]))
        ninety9s.append(float(content[6].strip().split(" ")[4]))
        errors.append(float(content[7].strip().split(" ")[3]))

    return [sum(throughputs), sum(medians) / numClients, sum(ninety9s) / numClients, sum(errors)]


def getQuePaxaPerformance(root, initClient, numClients):
    throughputs = []
    medians = []
    ninety9s = []
    errors = []
    for cl in list(range(initClient, initClient+numClients, 1)):
        file_name = root + str(cl) + ".log"
        # print(file_name + "\n")
        try:
            f = open(file_name, 'r')
        except OSError:
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        with f:
            content = f.readlines()
        if len(content) < 10:
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]
        if not (content[6].strip().split(" ")[0] == "Throughput" and content[7].strip().split(" ")[0] == "Median"):
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        throughputs.append(float(content[6].strip().split(" ")[2]))
        medians.append(float(content[7].strip().split(" ")[3]))
        ninety9s.append(float(content[8].strip().split(" ")[4]))
        errors.append(float(content[9].strip().split(" ")[3]))

    return [sum(throughputs), sum(medians) / numClients, sum(ninety9s) / numClients, sum(errors)]


def getEPaxosPaxosPerformance(root, initClient, numClients):
    throughputs = []
    medians = []
    ninety9s = []
    errors = []
    for cl in list(range(initClient, initClient+numClients, 1)):
        file_name = root + str(cl) + ".log"
        # print(file_name + "\n")
        try:
            f = open(file_name, 'r')
        except OSError:
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        with f:
            content = f.readlines()
        if len(content) < 5:
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]
        if not (content[1].strip().split(" ")[0] == "Throughput" and content[2].strip().split(" ")[0] == "Median"):
            print("Error in " + file_name + "\n")
            return [-1, -1, -1, -1]
        throughputs.append(float(content[1].strip().split(" ")[5]))
        medians.append(float(content[2].strip().split(" ")[3]))
        ninety9s.append(float(content[3].strip().split(" ")[4]))
        errors.append(float(content[4].strip().split(" ")[3]))

    return [sum(throughputs), sum(medians) / numClients, sum(ninety9s) / numClients, sum(errors)]