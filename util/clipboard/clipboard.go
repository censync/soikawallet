package clipboard

import (
	"errors"
	"os/exec"
)

var (
	selectedTool   = -1
	clipboardTools = [][]string{
		{"xsel", "--input", "--clipboard"},
		{"xclip", "-in", "-selection", "clipboard"},
		{"wl-copy"},
		{"termux-clipboard-set"},
	}
)

func init() {
	for i := 0; i < len(clipboardTools); i++ {
		if _, err := exec.LookPath(clipboardTools[0][0]); err == nil {
			selectedTool = i
			return
		}
	}
}

func IsClipBoardAvailable() bool {
	return selectedTool != -1
}

func CopyToClipboard(str string) error {
	if selectedTool == -1 {
		return errors.New(`cannot find clipboard tool, please, install "xsel", "xclip", "wl-copy" or "termux-clipboard-set" package`)
	}
	copyCmd := exec.Command(clipboardTools[selectedTool][0], clipboardTools[selectedTool][1:]...)
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
