package console_qr

import (
	"bytes"
	"github.com/skip2/go-qrcode"
	"io"
)

const (
	wb = "▄"
	ww = " "
	bw = "▀"
	bb = "█"

	frameSize = 192
)

func NewQR(str string) (string, error) {
	qrc, err := qrcode.New(str, qrcode.Low)
	if err != nil {
		return "", err
	}
	bitmap := qrc.Bitmap()

	buf := bytes.NewBufferString("")
	generateFrame(buf, bitmap)
	return buf.String(), nil
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
