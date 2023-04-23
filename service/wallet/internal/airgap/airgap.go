package airgap

const (
	ProtocolVersion = 1
)

type AirGap struct {
	protocol   uint8
	instanceId [33]byte
}

func Create() *AirGap {
	return &AirGap{
		protocol:   ProtocolVersion,
		instanceId: [33]byte{},
	}
}

func Restore(instanceId []byte) *AirGap {
	if len(instanceId) != 33 {
		panic("incorrect instance pub key size")
	}
	compressedPubKey := [33]byte{}
	copy(compressedPubKey[:], instanceId)
	return &AirGap{
		protocol:   ProtocolVersion,
		instanceId: compressedPubKey,
	}
}

type Message struct {
	Version    uint8
	InstanceId [33]byte
	Payload    []*OperationPayload
}

type OperationPayload struct {
	Operation uint8
	Size      uint32
	Data      []byte
}

func (a *AirGap) CreateMessage() *Message {
	return &Message{
		Version:    a.protocol,
		InstanceId: a.instanceId,
	}
}

func (m *Message) AddOperation(operation uint8, data []byte) *Message {
	m.Payload = append(m.Payload, &OperationPayload{
		Operation: operation,
		Size:      uint32(len(data)),
		Data:      data,
	})
	return m
}

func (m *Message) Bytes() []byte {
	result := make([]byte, 0)
	result = append(result, m.Version)
	result = append(result, m.InstanceId[:]...)
	for i := range m.Payload {
		payload := make([]byte, m.Payload[i].Size+1+4)

		payload[0] = m.Payload[i].Operation

		payload[1] = byte(m.Payload[i].Size >> 24)
		payload[2] = byte(m.Payload[i].Size >> 16)
		payload[3] = byte(m.Payload[i].Size >> 8)
		payload[4] = byte(m.Payload[i].Size)

		copy(payload[5:], m.Payload[i].Data)
		result = append(result, payload...)
	}
	return result
}
