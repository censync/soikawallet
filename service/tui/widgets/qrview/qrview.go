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
	*tview.TextView

	// vars
	redraw   func()
	chunks   []string
	duration time.Duration

	isStarted bool
	animDone  chan bool
}

func NewQrView(chunks []string, duration time.Duration, redraw func()) *QrView {
	layout := tview.NewTextView().
		SetWordWrap(false).
		//SetTextColor(tcell.ColorBlack).
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter)

	layout.SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorLightGray)

	return &QrView{
		TextView: layout,
		chunks:   chunks,
		duration: duration,
		redraw:   redraw,
	}

}

func NewQrViewText(data string) string {
	qrc, err := qrcode.New(data, qrcode.Low)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	bitmap := qrc.Bitmap()

	buf := bytes.NewBufferString("")
	generateFrame(buf, bitmap)
	return buf.String()
}

func (qr *QrView) Start() {
	if !qr.isStarted {
		qr.isStarted = true
	} else {
		return
	}

	qr.animDone = make(chan bool, 1)
	ticker := time.NewTicker(qr.duration * time.Millisecond)

	for i := 0; i < len(qr.chunks); i++ {
		qrc, err := qrcode.New(qr.chunks[i], qrcode.Low)
		if err != nil {
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		}
		bitmap := qrc.Bitmap()

		buf := bytes.NewBufferString("")
		generateFrame(buf, bitmap)
		qr.chunks[i] = buf.String()
	}

	if len(qr.chunks) > 1 {
		// show animation
		go func() {
			frameId := 0
			for {
				select {
				case <-ticker.C:
					frame := qr.chunks[frameId%len(qr.chunks)]
					qr.SetText(frame)
					qr.redraw()
					frameId++
				case <-qr.animDone:
					qr.isStarted = false
					ticker.Stop()
					close(qr.animDone)
					return
				}
			}
		}()
	} else {
		if len(qr.chunks) == 1 {
			qr.SetText(qr.chunks[0])
			qr.redraw()
		}
	}

}

func (qr *QrView) Stop() {
	qr.animDone <- true
	qr.Clear()
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
				_, _ = w.Write([]byte(bb))
			} else if currBlack && !nextBlack {
				_, _ = w.Write([]byte(bw))
			} else if !currBlack && !nextBlack {
				_, _ = w.Write([]byte(ww))
			} else {
				_, _ = w.Write([]byte(wb))
			}
		}
		_, _ = w.Write([]byte("\n"))
	}
}
