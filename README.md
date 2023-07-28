# QuePaxa

```QuePaxa``` is a novel crash fault tolerant and asynchronous consensus algorithm.

The main innovations of ```QuePaxa```, compared to the existing consensus algorithms such as ```Multi-Paxos```, ```Raft```, ```Rabia``` and ```EPaxos``` are threefold:

- ```QuePaxa``` employs a new asynchronous consensus core to tolerate adverse network conditions, while having a one-round-trip fast path under normal-case network conditions.

- ```QuePaxa``` employs ```hedging-delay``` instead of traditional ```timeouts```. ```Hedging delays``` can be arbitrarily configured in ```QuePaxa```, without affecting the liveness,
whereas for ```Raft``` and ```Multi-Paxos```, a conservatively high timeout should be used to ensure liveness.

- ```QuePaxa``` dynamically tunes the protocol at runtime to maximize the performance.


Our SOSP paper, ```QuePaxa: Escaping the tyranny of timeouts in consensus``` describes QuePaxa's design and evaluation in detail.

## Project Keywords

```state-machine replication (SMR)```,``` consensus```, ```asynchrony```, ```randomization```, ```tuning```, and ```hedging```

## Repository structure

- ```client/```, ```common/```, ```configuration/```, ```proto/```, ```replica/``` provide the ```QuePaxa``` ```replica``` and ```submitter``` implementations in Go-lang.

- ```integration-test/``` contains the set of integration tests that test the correctness of ```QuePaxa``` in a single machine deployment with 5 replicas and 5 clients under different
configuration parameters.

- ```experiments/```: contains the automated scripts for artifact evaluation and the compiled binaries of ```Epaxos```, ```Multi-paxos```, ```Raft```, ```Rabia``` and ```QuePaxa```.

## Documentation

- how to build, install and run ```QuePaxa``` in a single VM

- how to read ```QuePaxa``` codebase

### How to build and run QuePaxa

In this section, we explain how to build, install and run ```QuePaxa``` in a single VM

#### Single Machine Setup

##### Requirements

- a VM with at least 4 cores, 8GB memory and 10GB HDD
- Ubuntu 20.04 LTS with ```sudo``` access
- ```python3``` installed with ```numpy``` and ```matplotlib```

#####  Setup

- install ```python pip3```, ```matplotlib``` and ```go 1.19.5```
  
  ```
  sudo apt update
  sudo apt install python3-pip; pip3 install matplotlib
  rm -rf /usr/local/go
  wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.19.5.linux-amd64.tar.gz
  export PATH=$PATH:/usr/local/go/bin
  ```

  Check if the installation is successful by issueing the command: ```go version```
  which should output ```go1.19.5 linux/amd64```.

  This repository uses [Protocol Buffers](https://developers.google.com/protocol-buffers/).
  Install the protoc compiler by following the [Protocol Buffers](https://developers.google.com/protocol-buffers/).  

- clone QuePaxa from github and build the project
  
  ```
  git clone https://github.com/dedis/quepaxa
  cd quepaxa
  /bin/bash build.sh
  ```
  
  ```build.sh``` will produce an error ```protoc: command not found```, if you have not installed the ```protoc``` compiler correctly. 

  NOTE: The outputs might show ```raxos``` and ```QuePaxa``` interchangeably. This is because, ```QuePaxa``` code was initially names ```Raxos``` and later renamed to ```QuePaxa```. 

- run 5 replicas and 5 submitters in the same VM with ```50k cmd/sec``` aggregate arrival rate
  
  ```/bin/bash integration-test/safety-test.sh 200000 0 0 1 5000 50 10000 100 0 0```

  If the test was successful, the final outputs should look like the following.

  In the ```logs/200000/0/0/1/5000/50/10000/100/0/0/``` folder, you can see the output of 5 replicas (1-5.log) and 5 submitters (21-25.log).
  Submitter logs contain the median latency, throughput and 99 percentile latency statistics.

  A sample client output at ```logs/200000/0/0/1/5000/50/10000/100/0/0/21.log```

  ```
  initialized client 21 with process id: 31567
  starting request client
  Finish sending requests
  Calculating stats for thread 0
  Total time := 60 seconds
  Throughput := 10026 requests per second
  Median Latency := 5201 micro seconds per request
  99 pecentile latency := 19635 micro seconds per request
  Error Rate := 0

  finishing request client
  ```

  A sample replica output at ```logs/200000/0/0/1/5000/50/10000/100/0/0/1.log```
  ```
  started QuePaxa Server
  Average number of steps per slot: 1.000000, total slots 10183, steps accumilated 10183
  ```

  The log consistency (SMR correctness) can be found in the file ```logs/200000/0/0/1/5000/50/10000/100/0/0/consensus.log```
  ```
  60150596 entries match
  0 entries miss match

  TEST PASS
  ```

- Run integration tests that checks the SMR correctness of ```QuePaxa``` under different configuration parameters

  ```python3 integration-test/python/integration-automation.py```

  The integration tests consists of 9 tests, each of which exercise ```QuePaxa``` under different leader modes, timeouts and network conditions, using 5 replicas and 5 submitters setup in a single VM.  
  In the ```logs/``` folder you will find the log files corresponding to each test.
  Each subdirectory of ```logs/``` is indexed as ```logs/leaderTimeout/serverMode/leaderMode/pipeline/batchTime/batchSize/arrivalRate/closeLoopWindow/requestPropagationTime/asynchronousSimulationTime/```
  Please cross-check with the ```integration-test/python/integration-automation.py``` file to see what parameters are used in each integration subtest, that will help you to locate the output sub-folder in the ```logs/```  
  

### How to read QuePaxa code

In the following section, we outline the QuePaxa implementation architecture and the package level comments
that will allow you to understand the ```QuePaxa``` code base.

#### QuePaxa architecture

![architecture.jpg](docs%2Farchitecture.jpg)

Above figure depicts the ```QuePaxa``` architecture. ```QuePaxa``` contains two separate binaries;
the client (submitter in the paper) and the replica.

##### QuePaxa client

```QuePaxa``` client is the program that generates a stream of client requests to be consumed by the replica nodes. 
QuePaxa client supports the following configuration parameters.

- arrivalRate int: Poisson arrival rate in requests per second (default 10000)
- batchSize int : client batch size (default 50)
- batchTime int: client batch time in micro seconds (default 50)
- config string: quepaxa configuration file which contains the ip:port of each client and replica
- debugLevel int: debug level, debug messages with equal or higher debugLevel will be printed on console
- debugOn: turn on/off debug, turn off when benchmarking
- keyLen int: length of key in client requests (default 8)
- logFilePath string: log file path (default "logs/")
- name int: name of the client (default 21)
- operationType int: Type of operation for a status request: 1 bootstrap server, 2: print log, and 4. average number of
  steps per slot  (default 1)
- requestType string: request type: [status , request] (default "status")
- testDuration int: test duration in seconds (default 60)
- valLen int: length of value in client requests (default 8)
- window int: number of outstanding client batches sent by the client, before receiving the response (default 1000)


```QuePaxa``` client connects with all the QuePaxa replicas using TCP connections.

##### QuePaxa replica

```QuePaxa``` replica implements the ```QuePaxa``` consensus algorithm and the state machine replication logic.
```QuePaxa``` replica supports the following configurations.

- batchSize int: replica batch size (default 50)
- batchTime int: replica batch time in micro seconds (default 5000)
- benchmark int: Benchmark: 0 for KV store and 1 for Redis
- config string: QuePaxa configuration file (default "```configuration/local/configuration.yml```)
- debugLevel int: debug level (default 1010)
- debugOn: true / false
- epochSize int: epoch size for MAB (default 100)
- isAsync: true / false to simulate consensus level asynchrony
- keyLen int: length of key (default 8)
- leaderMode int: mode of leader change: 0 for fixed leader order, 1 for round robin, static partition,
  2 for M.A.B based on commit times, 3 for asynchronous, 4 for last committed proposer
- leaderTimeout int: leader timeout in micro seconds (default 5000000)
- logFilePath string: log file path (default "logs/")
- name int: name of the replica (default 1)
- pipelineLength int: pipeline length maximum number of outstanding proposals (default 1)
- requestPropagationTime int: additional wait time in 'milli seconds' for client batches,
  such that there is enough time for client driven request propagation
- serverMode int: 0 for non-lan-optimized, 1 for lan optimized
- timeEpochSize int: duration of a time epoch for the attacker in milli seconds (default 500)
- valLen int: length of value (default 8)
- asyncTimeOut int: artificial async timeout in milli seconds (default 500)

```QuePaxa``` replica consists of three layers: ```Proxy```, ```Proposer``` and ```Recorder```.

- ```Proxy```: Proxy maintains the state machine replication logic.
  Proxy receives client commands from the clients using TCP connections and forms a replica batch,
  to be consumed by the proposers.  
  Proxy then sends the replica batches to the proposer.  
  Upon reaching consensus, Proposer sends back the agreed upon value to the Proxy  
  and the Proxy updates the state machine. Then the proxy sends back the response to the client.

- ```Proposer```: Proposer is the implementation of the proposer segment of ```QuePaxa```,
  as described in our SOSP paper: QuePaxa: Escaping the tyranny of timeouts in consensus.
  Proposer communicates with all the recorders of all replicas using ```gRPC```.

- ```Recorder```: Recorder is the implementation of the ```Interval summary register``` as described in the our paper.
  Recorder responds to proposer requests.

#### Package comments

- ```client/```: contains the client implementation.
  The request sending logic is in ```request.go``` and the statistic calculation logic is in ```stat.go```

- ```common/```: contains go structs that are both common to replica and client

- ```configuration/```: contains the code to extract configuration data from a ```.yaml``` file into a
  ```cfg``` object.

- ```experiments/```: contains the AWS deployment scripts for the artifact evaluation.
   More on this is available in the artifact evaluation document.

- ```integration-test/```: contains the integration tests for ```QuePaxa```.

- ```proto/```: contains the proto definitions of client messages

- ```replica/```: contains the ```QuePaxa``` replica logic. ```proposer.go``` and ```recorder.go``` contain the
  proposer and the ```ISR``` logics, respectively.

- ```build.sh/```: contains the build script