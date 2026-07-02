package sysinfo

import (
	"os"
	"os/exec"
	"strings"
)

// GetSystemEditorPath gets editor using EDITOR env variable
// with common fallbacks if not set. The result may contain arguments
// (e.g. "code --wait"); callers should split on whitespace.
func GetSystemEditorPath() string {
	editor := strings.TrimSpace(os.Getenv("EDITOR"))
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
