import os


def base_case():
    print("testing base case with no batching and no pipelining, no view changes, fixed leader")
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 0
    pipeline = 1
    batchTime = 50
    batchSize = 1
    arrivalRate = 50
    closeLoopWindow = 2
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def fixed_leader():
    print("testing with batching and no pipelining, no view changes, fixed leader")
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 0
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def rotating_leader():
    print("testing with batching and no pipelining, no view changes, rotating leader")
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 1
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def mab_leader():
    print("testing with batching and no pipelining, no view changes, MAB leader")
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 2
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def last_proposed_leader():
    print("testing with batching and no pipelining, no view changes, last proposed leader")
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 4
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def asynchronous_leader():
    print("testing with batching and no pipelining, no view changes, no leaders")
    leaderTimeout = 1000
    serverMode = 0
    leaderMode = 3
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 2000
    closeLoopWindow = 50
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def optimization():
    print("testing with batching and no pipelining, no view changes, with the optimization")
    leaderTimeout = 3000000
    serverMode = 1
    leaderMode = 0
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 1
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def pipelining():
    print("testing with batching and pipelining, no view changes")
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 0
    pipeline = 10
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


def hedging():
    print("testing with batching and no pipelining, with hedging")
    leaderTimeout = 2000
    serverMode = 0
    leaderMode = 0
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime))


base_case()
fixed_leader()
rotating_leader()
mab_leader()
last_proposed_leader()
asynchronous_leader()
optimization()
pipelining()
hedging()
