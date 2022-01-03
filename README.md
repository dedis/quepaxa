Raxos is an asynchronous consensus protocol.

Raxos uses the [Mage](https://magefile.org/) build tool. Therefore it needs the ```mage``` command to be installed.


Raxos uses [Protocol Buffers](https://developers.google.com/protocol-buffers/).
It requires the ```protoc``` compiler with the ```go``` output plugin installed.


Raxos uses [Redis](https://redis.io/topics/quickstart) and it should be installed with default options.

All implementations are tested in ```Ubuntu 20.04.3 LTS```

Some build dependencies can be installed by running ```mage builddeps```

Download code dependencies by running ```mage deps``` or ```go mod vendor```.

Build the code using ```mage generate && mage build``` in the root directory

All the commands to run the server and the client are available in the respective directories