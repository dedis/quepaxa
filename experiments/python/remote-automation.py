import os

os.system("/bin/bash /home/pasindu/Documents/Raxos/experiments/remote/setup.sh")
print("Setup.sh completed\n")

for arrivalRate in [100, 200, 500, 600, 700, 800, 1000, 1200, 1500, 1800, 2000, 2200, 2500, 3000, 3500, 4000, 4500, 5000, 5500, 6000, 6500, 7000]:
    os.system("/bin/bash /home/pasindu/Documents/Raxos/experiments/remote/remote-test-overlay.sh " + str(arrivalRate))
    print(str(arrivalRate) + " completed\n")

print("The test completed")
