// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package tasks

import (
	"time"
)

const (
	Scheduler     = ConditionType(1)
	ChainEvent    = ConditionType(2)
	ChainGasPrice = ConditionType(3)

	ChainTransfer = OperationType(1)
	ChainCall     = OperationType(2)
	MultiOp       = OperationType(3)

	Active    = Status(1)
	Pending   = Status(2)
	Scheduled = Status(3)
	Retry     = Status(4)
	Completed = Status(5)
	Canceled  = Status(6)
)

type ConditionType uint8

type OperationType uint8

type Status uint8

type Task struct {
	Id        string        // unique task id
	Parent    string        // parent task id for multitasks
	Condition ConditionType // task condition for execution
	Operation OperationType // task operation
	Status    Status        // current status of task
	Retries   int           // -1 - unlimited
	Payload   interface{}
	Time      *time.Time
}

func (t *Task) hasFilter() bool {
	var result bool
	switch t.Operation {
	case ChainTransfer:
		result = true
	default:
		result = false
	}
	return result
}
