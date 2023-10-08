// Copyright 2023 The soikawallet Authors
// This file is part of soikawallet.
//
// soikawallet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// soikawallet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with  soikawallet. If not, see <http://www.gnu.org/licenses/>.

package events

import "fmt"

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
