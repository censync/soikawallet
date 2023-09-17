package event_bus

import "fmt"

const (
	EventLog EventType = iota + 100
	EventLogInfo
	EventLogSuccess
	EventLogWarning
	EventLogError

	EventUpdateCurrencies

	EventWalletInitialized
	EventDrawForce
	EventShowModal
	EventQuit

	// browser connector

	EventW3InternalConnections EventType = iota + 170
	EventW3Connect
	EventW3RequestAccounts
	EventW3ReqCallGetBlockByNumber

	EventW3Response
)

const (
	EventW3InternalGetConnections           = EventType(190)
	EventW3WalletAvailable        EventType = iota + 200
	EventW3WalletNotAvailable
	EventW3ConnAccepted
	EventW3ConnRejected
	EventW3CallGetBlockByNumber
)

type EventType uint8

type Event struct {
	event EventType
	data  interface{}
}

func (e *Event) Type() EventType {
	return e.event
}

func (e *Event) Data() interface{} {
	return e.data
}

func (e *Event) String() string {
	result, ok := e.data.(string)
	if !ok {
		result = fmt.Sprintf("%v", e.data)
	}
	return result
}

type EventBus struct {
	events  chan *Event
	done    chan bool
	stopped bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		events:  make(chan *Event, 100),
		stopped: false,
	}
}

func (b *EventBus) Stop() {
	if !b.stopped {
		b.stopped = true
		close(b.events)
	}
}

func (b *EventBus) IsStopped() bool {
	return b.stopped
}

func (b *EventBus) Events() <-chan *Event {
	return b.events
}

func (b *EventBus) Emit(event EventType, data interface{}) {
	if !b.stopped {
		b.events <- &Event{
			event: event,
			data:  data,
		}
	}
}
