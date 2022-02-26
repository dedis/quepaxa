import os

os.system("/bin/bash /home/pasindu/Documents/Raxos/experiments/remote/setup.sh")
print("Setup.sh completed\n")

for arrivalRate in [1000, 5000, 10000, 20000, 40000, 60000, 80000, 100000, 120000, 130000]:
    os.system("/bin/bash /home/pasindu/Documents/Raxos/experiments/remote/remote-test-overlay.sh " + str(arrivalRate))
    print(str(arrivalRate) + " completed\n")

print("The test complete")
