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

package clipboard

import (
	"errors"
	"io"
	"os/exec"
)

var (
	selectedCopyTool  = -1
	selectedPasteTool = -1
	clipboardCopy     = [][]string{
		{"xsel", "--input", "--clipboard"},
		{"xclip", "-in", "-selection", "clipboard"},
		{"wl-copy"},
		{"termux-clipboard-set"},
	}
	clipboardPaste = [][]string{
		//{"xsel", "--output", "--clipboard"},
		{"xclip", "-o"},
		{"wl-paste"},
		{"termux-clipboard-get"},
	}
)

func init() {
	for i := 0; i < len(clipboardCopy); i++ {
		if _, err := exec.LookPath(clipboardCopy[0][0]); err == nil {
			selectedCopyTool = i
			return
		}
	}
	for i := 0; i < len(clipboardPaste); i++ {
		if _, err := exec.LookPath(clipboardPaste[0][0]); err == nil {
			selectedPasteTool = i
			return
		}
	}
}

func Clear() {
	_ = CopyToClipboard(``)
}

func CopyToClipboard(str string) error {
	if selectedCopyTool == -1 {
		return errors.New(`cannot find clipboard tool, please, install "xsel", "xclip", "wl-copy" or "termux-clipboard-set" package`)
	}
	copyCmd := exec.Command(clipboardCopy[selectedCopyTool][0], clipboardCopy[selectedCopyTool][1:]...)
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = copyCmd.Start(); err != nil {
		return err
	}
	if _, err = in.Write([]byte(str)); err != nil {
		return err
	}
	if err = in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}

func PasteFromClipboard() (string, error) {
	if selectedPasteTool == -1 {
		return ``, errors.New(`cannot find clipboard tool, please, install "xsel", "xclip", "wl-paste" or "termux-clipboard-get" package`)
	}
	pasteCmd := exec.Command(clipboardCopy[selectedPasteTool][0], clipboardCopy[selectedPasteTool][1:]...)
	out, err := pasteCmd.StdoutPipe()
	if err != nil {
		return ``, err
	}

	defer out.Close()

	result, err := io.ReadAll(out)
	if err != nil {
		return ``, err
	}
	_ = pasteCmd.Wait()

	return string(result), nil
}
