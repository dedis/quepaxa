package proto

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

// MessageBlockRequest

func (t *MessageBlockRequest) Marshal(wire io.Writer) {

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

func (t *MessageBlockRequest) Unmarshal(wire io.Reader) error {

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
func (t *MessageBlockRequest) New() Serializable {
	return new(MessageBlockRequest)
}

// MessageBlock

func (t *MessageBlock) Marshal(wire io.Writer) {

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

func (t *MessageBlock) Unmarshal(wire io.Reader) error {

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
func (t *MessageBlock) New() Serializable {
	return new(MessageBlock)
}

// MessageBlockAck

func (t *MessageBlockAck) Marshal(wire io.Writer) {

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

func (t *MessageBlockAck) Unmarshal(wire io.Reader) error {

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
func (t *MessageBlockAck) New() Serializable {
	return new(MessageBlockAck)
}
