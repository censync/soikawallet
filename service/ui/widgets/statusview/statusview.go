package statusview

import (
	"fmt"
	"github.com/rivo/tview"
	"time"
)

const dateLogFormat = "15:04:05"

type StatusView struct {
	*tview.TextView
}

func NewStatusView() *StatusView {
	return &StatusView{TextView: tview.NewTextView()}
}

func (v *StatusView) Log(str string) {
	_, _ = fmt.Fprintf(v, "\n[white]%s:[white] [white]%s[white]", time.Now().Format(dateLogFormat), str)
}

func (v *StatusView) Info(str string) {
	_, _ = fmt.Fprintf(v, "\n[white]%s:[white] [blue]%s[white]", time.Now().Format(dateLogFormat), str)
}

func (v *StatusView) Success(str string) {
	_, _ = fmt.Fprintf(v, "\n[white]%s:[white] [green]%s[white]", time.Now().Format(dateLogFormat), str)
}

func (v *StatusView) Warn(str string) {
	_, _ = fmt.Fprintf(v, "\n[white]%s:[white] [orange]%s[white]", time.Now().Format(dateLogFormat), str)
}

func (v *StatusView) Error(str string) {
	_, _ = fmt.Fprintf(v, "\n[white]%s:[white] [red]%s[white]", time.Now().Format(dateLogFormat), str)
}
