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
        value = l.split(":")[1]
        if value != "no-op":
            dict[key] = value
            numberOfRequests = numberOfRequests + 1

    files.append(dict)
    print("Length of " + sys.argv[i] + " is " + str(numberOfRequests))
    print("Approximate throughput " + str(
        numberOfRequests / 60.0) + "requests per second")  # assuming a test duration of 60 todo change this if the test duration is changed


def equals(value1, value2):
    if value1 != value2:
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
                            print("mismatch in log position " + str(key) + " in " + mapName + ":" + map[
                                key] + " and " + tarName + ":" + files[j][key])
                            misMatch = misMatch + 1
                        else:
                            match = match + 1

    print(str(match) + " entries match")
    print(str(misMatch) + " entries miss match")


checkMaps(files)
