package proto

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

// ClientBatch

func (t *ClientBatch) Marshal(wire io.Writer) error {

	data, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	lengthWritten := len(data)
	var b [8]byte
	bs := b[:8]
	binary.LittleEndian.PutUint64(bs, uint64(lengthWritten))
	_, err = wire.Write(bs)
	if err != nil {
		return err
	}
	_, err = wire.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (t *ClientBatch) Unmarshal(wire io.Reader) error {

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
func (t *ClientBatch) New() Serializable {
	return new(ClientBatch)
}

// ClientStatus

func (t *ClientStatus) Marshal(wire io.Writer) error {

	data, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	lengthWritten := len(data)
	var b [8]byte
	bs := b[:8]
	binary.LittleEndian.PutUint64(bs, uint64(lengthWritten))
	_, err = wire.Write(bs)

	if err != nil {
		return err
	}
	_, err = wire.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (t *ClientStatus) Unmarshal(wire io.Reader) error {

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
func (t *ClientStatus) New() Serializable {
	return new(ClientStatus)
}
