import atexit
import os
import random
import sys
import time
from threading import Barrier
from threading import Thread

replica1 = sys.argv[1]
replica2 = sys.argv[2]
replica3 = sys.argv[3]
replica4 = sys.argv[4]
replica5 = sys.argv[5]

cert = sys.argv[6]
device_name = sys.argv[7]

reset = ["sshpass ssh -i " + cert + " " + replica1 + "  \"sudo tc qdisc del dev " + device_name + " root\"",
         "sshpass ssh -i " + cert + " " + replica2 + "  \"sudo tc qdisc del dev " + device_name + " root\"",
         "sshpass ssh -i " + cert + " " + replica3 + "  \"sudo tc qdisc del dev " + device_name + " root\"",
         "sshpass ssh -i " + cert + " " + replica4 + "  \"sudo tc qdisc del dev " + device_name + " root\"",
         "sshpass ssh -i " + cert + " " + replica5 + "  \"sudo tc qdisc del dev " + device_name + " root\""]

anormal = [
    "sshpass ssh -i  " + cert + " " + replica1 + "  \"sudo tc qdisc add dev " + device_name + " root netem delay 500ms\"",
    "sshpass ssh -i  " + cert + " " + replica2 + "  \"sudo tc qdisc add dev " + device_name + " root netem delay 500ms\"",
    "sshpass ssh -i  " + cert + " " + replica3 + "  \"sudo tc qdisc add dev " + device_name + " root netem delay 500ms\"",
    "sshpass ssh -i  " + cert + " " + replica4 + "  \"sudo tc qdisc add dev " + device_name + " root netem delay 500ms\"",
    "sshpass ssh -i  " + cert + " " + replica5 + "  \"sudo tc qdisc add dev " + device_name + " root netem delay 500ms\""]


def exit_handler():
    sys.stdout.flush()
    for i in range(len(reset)):
        os.system(reset[i])


def execute(cmd, barrier):
    os.system(cmd)
    barrier.wait()


atexit.register(exit_handler)

sys.stdout.flush()
exit_handler()
t_end = time.time() + 30
while time.time() < t_end:
    randomInstance1 = random.randint(0, 4)
    randomInstance2 = random.randint(0, 4)
    while randomInstance1 == randomInstance2:
        randomInstance2 = random.randint(0, 4)

    barrier = Barrier(2 + 1)

    worker1 = Thread(target=execute, args=(anormal[randomInstance1], barrier))
    worker2 = Thread(target=execute, args=(anormal[randomInstance2], barrier))

    worker1.start()
    worker2.start()

    barrier.wait()

    barrier = Barrier(2 + 1)

    worker1 = Thread(target=execute, args=(reset[randomInstance1], barrier))
    worker2 = Thread(target=execute, args=(reset[randomInstance2], barrier))

    worker1.start()
    worker2.start()

    barrier.wait()
