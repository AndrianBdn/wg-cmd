package sysinfo

import (
	"os"
	"os/exec"
)

func IsRoot() bool {
	return os.Getuid() == 0
}

func binPath(binary string) string {
	path, err := exec.LookPath(binary)
	if err == nil {
		return path
	}
	candidates := []string{"/usr/bin/" + binary, "/bin/" + binary, "/usr/sbin/" + binary, "/sbin/" + binary}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func HasWireguard() bool {
	if os.Getenv("WG_CMD_NO_DEPS") != "" {
		return true
	}
	return binPath("wg") != ""
}
