import sys

files = []
for i in range(1, len(sys.argv)):
    dict = {}
    lines = []
    with open(sys.argv[i]) as file:
        lines = [line.rstrip() for line in file]

    numberOfRequests = 0
    for l in lines:
        key = l.split(":")[0]
        values = l.split(":")[1].split(",")
        dict[key] = values[:len(values) - 1]
        numberOfRequests = numberOfRequests + len(values) - 1

    files.append(dict)
    print("Length of " + sys.argv[i] + " is " + str(numberOfRequests))
    print("Approximate throughput " + str(numberOfRequests / 60.0) + "requests per second")


def equals(array1, array2):
    if len(array1) != len(array2):
        return False
    for i in range(len(array1)):
        if array1[i] != array2[i]:
            return False

    return True


def checkMaps(files):
    misMatch = 0
    match = 0
    for i in range(len(files)):
        map = files[i]
        mapName = sys.argv[i + 1]
        for key in map.keys():
            for j in range(len(files)):
                if i == j:
                    continue
                else:
                    tarName = sys.argv[j + 1]
                    if key in files[j].keys():
                        if not equals(files[j][key], map[key]):
                            print("Mismatch in log position " + str(key) + " in " + mapName + ":" + map[
                                key] + " and " + tarName + ":" + files[j][key])
                            misMatch = misMatch + len(map[key])
                        else:
                            match = match + len(map[key])

    print(str(misMatch) + " entries miss match")


checkMaps(files)
