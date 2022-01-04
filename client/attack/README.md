Attacker client directly manipulates the network interface card.

To run attack experiment, the consensus replicas should be run in three different virtual machines, and the attacker should be run in a seperate virtual machine.

```.cmd/attackerctl/bin/attackctl send "-" --config doc/configuration/local/backsosctl/quorum.yml --attackType delay ```