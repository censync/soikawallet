package handler

import "fmt"

const (
	EventLog EventType = iota
	EventLogInfo
	EventLogSuccess
	EventLogWarning
	EventLogError

	EventWalletInitialized
	EventDrawForce
	EventShowModal
	EventQuit
)

type EventType uint8

type TuiEvent struct {
	event EventType
	data  interface{}
}

type TBus chan *TuiEvent

func (t *TuiEvent) Type() EventType {
	return t.event
}

func (t *TuiEvent) Data() interface{} {
	return t.data
}

func (t *TuiEvent) String() string {
	result, ok := t.data.(string)
	if !ok {
		result = fmt.Sprintf("%v", t.data)
	}
	return result
}

func (b *TBus) Emit(event EventType, data interface{}) {
	*b <- &TuiEvent{
		event: event,
		data:  data,
	}
}
