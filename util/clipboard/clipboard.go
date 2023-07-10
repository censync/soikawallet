package clipboard

import (
	"errors"
	"os/exec"
)

var (
	selectedTool  = -1
	clipboardCopy = [][]string{
		{"xsel", "--input", "--clipboard"},
		{"xclip", "-in", "-selection", "clipboard"},
		{"wl-copy"},
		{"termux-clipboard-set"},
	}
	clipboardPaste = [][]string{
		{"xsel", "--output", "--clipboard"},
		{"xclip", "-o"},
		{"wl-paste"},
		{"termux-clipboard-get"},
	}
)

func init() {
	for i := 0; i < len(clipboardCopy); i++ {
		if _, err := exec.LookPath(clipboardCopy[0][0]); err == nil {
			selectedTool = i
			return
		}
	}
}

func Clear() {
	CopyToClipboard(``)
}

func CopyToClipboard(str string) error {
	if selectedTool == -1 {
		return errors.New(`cannot find clipboard tool, please, install "xsel", "xclip", "wl-copy" or "termux-clipboard-set" package`)
	}
	copyCmd := exec.Command(clipboardCopy[selectedTool][0], clipboardCopy[selectedTool][1:]...)
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
