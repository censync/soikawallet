package qrview

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
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

func ShowAnimation(str *string, duration int, redraw func()) *tview.TextView {

	labelQr := tview.NewTextView().
		SetWordWrap(false).
		//SetTextColor(tcell.ColorBlack).
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter)

	labelQr.SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorLightGray)

	update := func() {
		for i := 0; i < 15; i++ {
			payload := make([]byte, frameSize)
			rand.Read(payload)
			b64 := base64.StdEncoding.EncodeToString(payload)
			qrc, err := qrcode.New(b64, qrcode.Low)
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

func generateAnimation(data []byte) [][][]bool {
	for iter := 0; iter < len(data); iter += frameSize {

	}
	if len(data)%frameSize > 0 {
		// append
	}
	if len(data) < frameSize {
		// append
	}
	return nil
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
