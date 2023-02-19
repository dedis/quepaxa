import os


def run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate):
    os.system("/bin/bash /home/pasindu/Documents/Raxos/experiment/local/local-test.sh "+ str(epoch) + " "
              + str(leaderTimeout) + " "
              + str(serverMode) + " "
              + str(leaderMode) + " "
              + str(pipeline) + " "
              + str(batchTime) + " "
              + str(batchSize) + " "
              + str(arrivalRate))



# case 1

epoch = 200


leaderTimeout = 3000000
serverMode = 0
leaderMode = 0
pipeline = 1
batchTime = 600
batchSize = 1
arrivalRate = 1000
# run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)

# case 2

leaderTimeout = 3000000
serverMode = 0
leaderMode = 0
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 20000
# run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)


# case 3

leaderTimeout = 3000000
serverMode = 0
leaderMode = 0
pipeline = 10
batchTime = 2000
batchSize = 50
arrivalRate = 20000
# run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)


# case 4

leaderTimeout = 200000
serverMode = 0
leaderMode = 1
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 20000
run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)


# case 5

leaderTimeout = 200000
serverMode = 0
leaderMode = 2
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 20000
run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)

# case 6

leaderTimeout = 200000
serverMode = 0
leaderMode = 3
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 10000
run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)


# case 7

leaderTimeout = 3000000
serverMode = 1
leaderMode = 0
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 20000
# run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)


# case 8

leaderTimeout = 30000
serverMode = 0
leaderMode = 0
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 20000
# run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)

# case 9

leaderTimeout = 3000
serverMode = 0
leaderMode = 0
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 20000
# run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)


# case 10

leaderTimeout = 300
serverMode = 0
leaderMode = 0
pipeline = 1
batchTime = 2000
batchSize = 50
arrivalRate = 20000
# run(epoch, leaderTimeout, serverMode, leaderMode, pipeline, batchTime, batchSize, arrivalRate)


print("Test completed")