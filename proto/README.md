This directory defines the protobuf structs used for client ```requests``` and client ```status``` messages.

If you modify the ```.proto``` files, run ```mage build``` to regenerate the protobuff structs.

For import modularity, we moved the replica-replica message types ```.proto``` files to ```replica/``` directory