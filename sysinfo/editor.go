package sysinfo

import (
	"os"
	"os/exec"
)

// GetSystemEditorPath gets editor using EDITOR env variable
// with common fallbacks if not set
func GetSystemEditorPath() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editors := []string{"vim", "vi", "nano", "pico", "ed", "emacs"}
		for _, e := range editors {
			// check if editor exists using exec.LookPath
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}
	return editor
}
