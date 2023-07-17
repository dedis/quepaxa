QuePaxa is an asynchronous consensus protocol.  

[//]: # (The corresponding paper is in the [quepaxa-doc]&#40;https://github.com/dedis/quepaxa-doc&#41; repository.)

QuePaxa uses the [Mage](https://magefile.org/) build tool. Therefore it needs the ```mage``` command to be installed.


QuePaxa uses [Protocol Buffers](https://developers.google.com/protocol-buffers/).
It requires the ```protoc``` compiler with the ```go``` output plugin installed.


QuePaxa uses [Redis](https://redis.io/topics/quickstart) and it should be installed with default options.

All implementations are tested in ```Ubuntu 20.04.3 LTS```

Some build dependencies can be installed by running ```mage builddeps```

Download code dependencies by running ```mage deps``` or ```go mod vendor```.

Build the code using ```mage generate && mage build``` in the root directory

All the commands to run the server and the client are available in the ```integration-test``` directory
