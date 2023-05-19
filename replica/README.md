This folder contains the QuePaxa replica code

QuePaxa implementation contains three components: (1) proxy, (2) proposer, and (3) recorder

Proxy is the component which receives client batches from clients, forms replica batches, and send to Proposer for replication

Proposer is the implementation of the exact proposer in the QuePaxa paper

Recorder is the implementation of the interval summary register (ISR) in the QuePaxa paper.

For instructions on how to run a QuePaxa replica, please refer the ```integration-test/safety-test.sh```