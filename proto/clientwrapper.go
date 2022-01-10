package proto

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

// ClientRequest

func (t *ClientRequestBatch) Marshal(wire io.Writer) {

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

func (t *ClientRequestBatch) Unmarshal(wire io.Reader) error {

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
func (t *ClientRequestBatch) New() Serializable {
	return new(ClientRequestBatch)
}

// ClientResponse

func (t *ClientResponseBatch) Marshal(wire io.Writer) {

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

func (t *ClientResponseBatch) Unmarshal(wire io.Reader) error {

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
func (t *ClientResponseBatch) New() Serializable {
	return new(ClientResponseBatch)
}
