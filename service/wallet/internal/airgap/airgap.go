package airgap

import (
	"bytes"
	"errors"
)

const (
	VersionDefault       = 1
	compressedPubKeySize = 33
)

type AirGap struct {
	// version of protocol
	version uint8
	// instanceId compressed public key for device pairing
	instanceId []byte

	ed EncryptorDecryptor
}

type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
}

type Decryptor interface {
	Decrypt(data []byte) ([]byte, error)
}

type EncryptorDecryptor interface {
	Encryptor
	Decryptor
}

type Message struct {
	Version    uint8
	InstanceId []byte
	Payload    []*OpPayload
	e          Encryptor
}

// OpPayload is operation payload data
type OpPayload struct {
	// OpCode - operation code
	OpCode uint16
	Size   uint32
	Data   []byte
}

// NewAirGap initiates a new AirGap instance with secp256k1 serialized compressed public key
func NewAirGap(version uint8, instanceId []byte) *AirGap {
	//var compressedPubKey []byte
	if instanceId == nil || len(instanceId) != compressedPubKeySize {
		panic("incorrect instance pub key size")
	}

	//copy(compressedPubKey, instanceId)
	return &AirGap{
		version:    version,
		instanceId: instanceId,
	}
}

func (a *AirGap) SetEncryptorDecryptor(ed EncryptorDecryptor) *AirGap {
	a.ed = ed
	return a
}

// CreateMessage initiates new builder for AirGap messages batch
func (a *AirGap) CreateMessage() *Message {
	if a.instanceId == nil {
		panic("instance id is not defined")
	}
	return &Message{
		Version:    a.version,
		InstanceId: a.instanceId,
		e:          a.ed,
	}
}

func (m *Message) AddOperation(opCode uint16, data []byte) *Message {
	m.Payload = append(m.Payload, &OpPayload{
		OpCode: opCode,
		Size:   uint32(len(data)),
		Data:   data,
	})
	return m
}

func (m *Message) Marshal() ([]byte, error) {
	result := make([]byte, 0)
	result = append(result, m.Version)
	result = append(result, m.InstanceId[:]...)
	for i := range m.Payload {
		payload := make([]byte, m.Payload[i].Size+2+4)

		payload[0] = byte(m.Payload[i].OpCode >> 8)
		payload[1] = byte(m.Payload[i].OpCode)

		payload[2] = byte(m.Payload[i].Size >> 24)
		payload[3] = byte(m.Payload[i].Size >> 16)
		payload[4] = byte(m.Payload[i].Size >> 8)
		payload[5] = byte(m.Payload[i].Size)

		copy(payload[6:], m.Payload[i].Data)
		result = append(result, payload...)
	}

	if m.e != nil {
		return m.e.Encrypt(result)
	}
	return result, nil
}

func (m *Message) MarshalB64Chunks() ([]string, error) {
	serializedMessages, err := m.Marshal()
	if err != nil {
		return nil, err
	}

	result, err := NewChunks(serializedMessages, defaultChunkSize)

	if err != nil {
		return nil, err
	}

	return result.SerializeB64(), nil
}

func (a *AirGap) Unmarshal(data []byte) (*Message, error) {
	var err error

	if a.ed != nil {
		data, err = a.ed.Decrypt(data)
		if err != nil {
			return nil, err
		}
	}
	version := data[0]
	instanceId := data[1 : compressedPubKeySize+1]

	if version != a.version {
		if version < a.version {
			return nil, errors.New("airgap message version less than supported")
		}

		if version > a.version {
			return nil, errors.New("airgap message version greater than supported")
		}
	}

	if !bytes.Equal(a.instanceId, instanceId) {
		return nil, errors.New("airgap message has incorrect instance")
	}
	message := a.CreateMessage()

	bytesReaded := compressedPubKeySize + 1

	for i := bytesReaded; i < len(data); i += bytesReaded {
		opCode := uint16(data[i+1]) | uint16(data[i])<<8
		size := uint32(data[i+5]) | uint32(data[i+4])<<8 | uint32(data[i+3])<<16 | uint32(data[i+2])<<24
		bytesReaded = 6 + int(size)
		message.AddOperation(opCode, data[i+6:i+bytesReaded])

	}

	return message, nil
}
