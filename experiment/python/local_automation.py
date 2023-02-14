import os

epoch = 20
leaderTimeout = 3000000
serverMode = 0
leaderMode = 0
pipeline = 1
batchTime = 1
batchSize = 1
arrivalRate = 1000

os.system("/bin/bash /home/pasindu/Documents/Raxos/experiment/local/local-test.sh " + str(epoch) + " " 
          + str(leaderTimeout) + " "
          + str(serverMode) + " "
          + str(leaderMode) + " "
          + str(pipeline) + " "
          + str(batchTime) + " "
          + str(batchSize) + " "
          + str(arrivalRate))

print("Test completed")
