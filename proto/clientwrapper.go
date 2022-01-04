package proto

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

// ClientRequest

func (t *ClientRequest) Marshal(wire io.Writer) {

	data, err := proto.Marshal(t)
	if err != nil {
		return
	}
	lengthWritten := len(data)
	var b [4]byte
	bs := b[:4]
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

func (t *ClientRequest) Unmarshal(wire io.Reader) error {

	var b [4]byte
	bs := b[:4]

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
func (t *ClientRequest) New() Serializable {
	return new(GenericConsensus)
}

// ClientResponse

func (t *ClientResponse) Marshal(wire io.Writer) {

	data, err := proto.Marshal(t)
	if err != nil {
		return
	}
	lengthWritten := len(data)
	var b [4]byte
	bs := b[:4]
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

func (t *ClientResponse) Unmarshal(wire io.Reader) error {

	var b [4]byte
	bs := b[:4]

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
func (t *ClientResponse) New() Serializable {
	return new(GenericConsensus)
}
