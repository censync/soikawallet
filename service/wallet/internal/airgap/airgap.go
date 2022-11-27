package airgap

import "github.com/google/uuid"

type AirGap struct {
	protocol   uint8
	instanceId [16]byte
}

func Create() *AirGap {
	id := uuid.New()
	return &AirGap{
		protocol:   ProtocolVersion,
		instanceId: id,
	}
}

func Restore(instanceId string) *AirGap {
	id, err := uuid.Parse(instanceId)
	if err != nil {
		panic("cannot restore instance id")
	}
	return &AirGap{
		protocol:   ProtocolVersion,
		instanceId: id,
	}
}

func (a *AirGap) CreateMessage(operation uint8, data []byte) *Message {
	return &Message{
		Version:    a.protocol,
		InstanceId: a.instanceId,
		Operation:  operation,
		Data:       data,
	}
}
