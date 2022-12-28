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
	spinnerDone chan struct{}
	isStarted   bool
}

func NewSpinner(spinnerType int, interval time.Duration) *Spinner {
	if spinnerType < 0 || spinnerType > len(frames)-1 {
		panic("incorrect spinner type ")
	}
	return &Spinner{spinnerType: spinnerType, interval: interval}
}

func (s *Spinner) Start(callback func(string)) {
	s.spinnerDone = make(chan struct{})
	ticker := time.NewTicker(s.interval * time.Millisecond)
	if !s.isStarted {
		s.isStarted = true
	} else {
		return
	}

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
	<-s.spinnerDone
}
