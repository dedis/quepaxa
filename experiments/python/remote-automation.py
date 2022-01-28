import os

os.system("/bin/bash /home/oem/GolandProjects/raxos/experiments/remote/setup.sh")
print("Setup.sh completed\n")
for arrivalRate in [1000, 5000, 10000, 15000, 25000, 55000, 75000, 100000, 110000, 120000, 130000, 140000, 150000, 170000, 180000, 190000, 200000]:
    os.system("/bin/bash /home/oem/GolandProjects/raxos/experiments/remote/remote-test-overlay.sh " + str(arrivalRate))
    print(str(arrivalRate) + " completed\n")