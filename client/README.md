This directory implements the QuePaxa client.


Client implementation supports two operations.

(1) Send a ```status``` request to proxies

(2) Send SMR ```request```s to proxies

```status``` client is used to send status operations: (1) server bootstrapping, (2) log printing, and (3) printing server side stats.

For instructions on how to run the client, please refer the ```test/correctness-test.sh```