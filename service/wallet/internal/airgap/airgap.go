package airgap

const (
	ProtocolVersion = 1

	OperationInit        = 1
	OperationHandshake   = 2
	OperationImportMeta  = 3
	OperationExportMeta  = 4
	OperationSignRequest = 5
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
	return &AirGap{
		protocol:   ProtocolVersion,
		instanceId: [33]byte{},
	}
}

type Message struct {
	Version    uint8
	InstanceId [33]byte
	Operation  uint8
	Data       []byte
}

func (a *AirGap) CreateMessage(operation uint8, data []byte) *Message {
	return &Message{
		Version:    a.protocol,
		InstanceId: a.instanceId,
		Operation:  operation,
		Data:       data,
	}
}

func (m *Message) Bytes() []byte {
	result := make([]byte, 0)
	result = append(result, m.Version)
	result = append(result, m.InstanceId[:]...)
	result = append(result, m.Operation)
	result = append(result, m.Data...)
	return result
}
