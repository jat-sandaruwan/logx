package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Open opens a file in the appropriate editor based on platform
func Open(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Try Notepad++ first, fall back to notepad
		notepadpp := `C:\Program Files\Notepad++\notepad++.exe`
		if _, err := os.Stat(notepadpp); err == nil {
			cmd = exec.Command(notepadpp, filePath)
		} else {
			notepadpp = `C:\Program Files (x86)\Notepad++\notepad++.exe`
			if _, err := os.Stat(notepadpp); err == nil {
				cmd = exec.Command(notepadpp, filePath)
			} else {
				cmd = exec.Command("notepad", filePath)
			}
		}
	case "darwin":
		// macOS - try VS Code, Sublime, TextEdit
		editors := []string{"code", "subl", "open"}
		for _, editor := range editors {
			if _, err := exec.LookPath(editor); err == nil {
				if editor == "open" {
					cmd = exec.Command(editor, "-e", filePath)
				} else {
					cmd = exec.Command(editor, filePath)
				}
				break
			}
		}
		if cmd == nil {
			cmd = exec.Command("open", "-e", filePath)
		}
	case "linux":
		// Linux - try various editors
		editors := []string{"code", "gedit", "kate", "nano", "vim"}
		for _, editor := range editors {
			if _, err := exec.LookPath(editor); err == nil {
				cmd = exec.Command(editor, filePath)
				break
			}
		}
		if cmd == nil {
			return fmt.Errorf("no suitable editor found")
		}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// OpenWithCustom opens a file with a custom editor command
func OpenWithCustom(filePath, editorCmd string) error {
	cmd := exec.Command(editorCmd, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}