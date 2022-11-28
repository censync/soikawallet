package qrview

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"time"
)

const (
	wb = "▄"
	ww = " "
	bw = "▀"
	bb = "█"

	frameSize = 192
)

type QrView struct {
	*tview.Flex
}

func ShowAnimation(chunks []string, duration int, redraw func()) *tview.TextView {

	labelQr := tview.NewTextView().
		SetWordWrap(false).
		//SetTextColor(tcell.ColorBlack).
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter)

	labelQr.SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorLightGray)

	update := func() {

		for i := 0; i < len(chunks); i++ {
			qrc, err := qrcode.New(chunks[i], qrcode.Low)
			if err != nil {
				log.Fatal(err)
			}
			bitmap := qrc.Bitmap()

			buf := bytes.NewBufferString("")
			generateFrame(buf, bitmap)
			labelQr.SetText(buf.String())
			redraw()
			time.Sleep(time.Duration(duration) * time.Millisecond)
		}
	}

	go update()

	return labelQr
}

func generateFrame(w io.Writer, code [][]bool) {
	size := len(code) - 1
	for i := 0; i <= size; i += 2 {
		for j := 0; j <= size; j++ {
			nextBlack := false
			if i+1 < size {
				nextBlack = code[j][i+1]
			}
			currBlack := code[j][i]
			if currBlack && nextBlack {
				w.Write([]byte(bb))
			} else if currBlack && !nextBlack {
				w.Write([]byte(bw))
			} else if !currBlack && !nextBlack {
				w.Write([]byte(ww))
			} else {
				w.Write([]byte(wb))
			}
		}
		w.Write([]byte("\n"))
	}
}
