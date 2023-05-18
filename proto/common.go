package proto

import "io"

/*
	Each Message (except for gRPC messages) sent over the network should implement this interface
	If a new message type needs to be added: first define it in a proto file, generate the go protobuf files using "mage generate" and then implement the following three methods

*/

type Serializable interface {
	Marshal(io.Writer) error
	Unmarshal(io.Reader) error
	New() Serializable
}
