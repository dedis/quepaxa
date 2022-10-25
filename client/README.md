Client implementation supports two operations.

(1) Send a ```status``` request to proxies

(2) Send SMR ```request```s to proxies

To send a status request ```./client/bin/client --name 11 --requestType status --operationType [1, 2]```

```OperationType 1``` for server bootstrapping 

```OperationType 2``` for server log printing

To send client requests with minimal options ```./client/bin/client --name 11 --requestType request```

Supported options

```--name```: name of the client as specified in the ```configuration.yml```

```--config```: configuration file

```--logFilePath```: log file path

```--batchSize```: client batch size 

```--testDuration```: test duration in seconds

```--arrivalRate```: poisson arrival rate in requests per second

```--requestType```: "request type: [```status``` , ```request```]

```--operationType```: Type of operation for a status request: ```1``` (bootstrap server), ```2```: (print log)

```--keyLen```: key length

```--valLen```: value length

```--debugOn```: ```false``` or ```true```

```--debugLevel```: debug level

 