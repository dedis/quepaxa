import os


def base_case():
    print("testing base case with no batching and no pipelining, no view changes, fixed leader", flush=True)
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 0
    pipeline = 1
    batchTime = 50
    batchSize = 1
    arrivalRate = 50
    closeLoopWindow = 2
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def fixed_leader():
    print("testing with batching and no pipelining, no view changes, fixed leader", flush=True)
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 0
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def rotating_leader():
    print("testing with batching and no pipelining, no view changes, rotating leader", flush=True)
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 1
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def mab_leader():
    print("testing with batching and no pipelining, no view changes, MAB leader", flush=True)
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 2
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def last_proposed_leader():
    print("testing with batching and no pipelining, no view changes, last proposed leader", flush=True)
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 4
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def asynchronous_leader():
    print("testing with batching and no pipelining, no view changes, no leaders", flush=True)
    leaderTimeout = 1000
    serverMode = 0
    leaderMode = 3
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 2000
    closeLoopWindow = 50
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def optimization():
    print("testing with batching and no pipelining, no view changes, with the optimization", flush=True)
    leaderTimeout = 3000000
    serverMode = 1
    leaderMode = 0
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 1
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def pipelining():
    print("testing with batching and pipelining, no view changes", flush=True)
    leaderTimeout = 3000000
    serverMode = 0
    leaderMode = 0
    pipeline = 10
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


def hedging():
    print("testing with batching and no pipelining, with hedging", flush=True)
    leaderTimeout = 2000
    serverMode = 0
    leaderMode = 0
    pipeline = 1
    batchTime = 1000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 100
    requestPropagationTime = 0
    asynchronousSimulationTime = 0
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))

def asynchronous():
    print("testing with batching and no pipelining, with asynchronous attack", flush=True)
    leaderTimeout = 2000
    serverMode = 0
    leaderMode = 4
    pipeline = 1
    batchTime = 5000
    batchSize = 50
    arrivalRate = 10000
    closeLoopWindow = 1000
    requestPropagationTime = 0
    asynchronousSimulationTime = 2
    os.system("/bin/bash integration-test/safety-test.sh " + str(leaderTimeout) + " " + str(serverMode) + " " + str(
        leaderMode) + " " + str(pipeline) + " " + str(batchTime) + " " + str(batchSize) + " " + str(arrivalRate) +
              " " + str(closeLoopWindow) + " " + str(requestPropagationTime)+ " "+str(asynchronousSimulationTime))


base_case()
fixed_leader()
rotating_leader()
mab_leader()
last_proposed_leader()
asynchronous_leader()
optimization()
pipelining()
hedging()
asynchronous()
