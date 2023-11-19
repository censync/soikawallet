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
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type TaskHub struct {
	tasks map[string]*Task
	mu    sync.RWMutex
}

var (
	log = logrus.WithFields(logrus.Fields{"service": "connector", "module": "tasks"})
)

func (t *TaskHub) Add(task *Task) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if _, ok := t.tasks[task.Id]; ok {
		return errors.New("task already exist")
	}

	t.tasks[task.Id] = task
	return nil
}
