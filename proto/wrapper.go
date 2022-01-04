package proto

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

type Serializable interface {
	Marshal(io.Writer)
	Unmarshal(io.Reader) error
	New() Serializable
}

// Prepare

func (t *consensus.GenericConsensus) Marshal(wire io.Writer) {
	data, err := proto.Marshal(t)
	if err != nil {
		return
	}
	lengthWritten := len(data)
	var b [8]byte
	bs := b[:8]
	binary.LittleEndian.PutUint64(bs, uint64(lengthWritten))
	_, err = wire.Write(bs)
	if err != nil {
		return
	}
	_, err = wire.Write(data)
	if err != nil {
		return
	}
}

func (t *GenericConsensus) Unmarshal(wire io.Reader) error {

	var b [8]byte
	bs := b[:8]

	_, err := io.ReadFull(wire, bs)
	if err != nil {
		return err
	}
	numBytes := binary.LittleEndian.Uint64(bs)

	data := make([]byte, numBytes)
	length, err := io.ReadFull(wire, data)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data[:length], t)
	if err != nil {
		return err
	}
	return nil
}
func (t *GenericConsensus) New() Serializable {
	return new(Prepare)
}
