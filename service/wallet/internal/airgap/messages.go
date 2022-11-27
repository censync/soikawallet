package airgap

type Message struct {
	Version    uint8
	InstanceId [16]byte
	Operation  uint8
	Data       []byte
}

func (m *Message) Bytes() []byte {
	result := make([]byte, 0)
	result = append(result, m.Version)
	result = append(result, m.InstanceId[:]...)
	result = append(result, m.Operation)
	result = append(result, m.Data...)
	return result
}
