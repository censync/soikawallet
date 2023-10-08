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
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package spinner

import (
	"time"
)

const (
	SpinOne = iota
	SpinThree
)

var frames = [][]string{
	{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	{".  ", ".. ", "...", " ..", "  .", "   "},
}

type Spinner struct {
	spinnerType int
	interval    time.Duration
	spinnerDone chan bool
	isStarted   bool
}

// TODO: Add timeout
func NewSpinner(spinnerType int, interval time.Duration) *Spinner {
	if spinnerType < 0 || spinnerType > len(frames)-1 {
		panic("incorrect spinner type ")
	}
	return &Spinner{spinnerType: spinnerType, interval: interval}
}

func (s *Spinner) Start(callback func(string)) {
	if !s.isStarted {
		s.isStarted = true
	} else {
		return
	}

	s.spinnerDone = make(chan bool, 1)
	ticker := time.NewTicker(s.interval * time.Millisecond)

	go func() {
		frameId := 0
		for {
			select {
			case <-ticker.C:
				frame := frames[s.spinnerType][frameId%len(frames[s.spinnerType])]
				callback(frame)
				frameId++
			case <-s.spinnerDone:
				s.isStarted = false
				ticker.Stop()
				close(s.spinnerDone)
				return
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.spinnerDone <- true
}
