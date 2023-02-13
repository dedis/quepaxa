Client implementation supports two operations.

(1) Send a ```status``` request to proxies

(2) Send SMR ```request```s to proxies

To send a status request ```./client/bin/client --name 11 --requestType status --operationType [1, 2]```

```OperationType 1``` for server bootstrapping 

```OperationType 2``` for server log printing

To send client requests with minimal options ```./client/bin/client --name 11 --requestType request```