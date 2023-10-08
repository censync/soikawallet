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

package statusview

import (
	"fmt"
	"github.com/censync/tview"
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
