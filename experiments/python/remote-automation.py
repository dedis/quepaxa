import os

os.system("/bin/bash /home/pasindu/Documents/Raxos/experiments/remote/setup.sh")
print("Setup.sh completed\n")

for arrivalRate in [120000, 130000]:
    os.system("/bin/bash /home/pasindu/Documents/Raxos/experiments/remote/remote-test-overlay.sh " + str(arrivalRate))
    print(str(arrivalRate) + " completed\n")

print("The test complete")
