package sysinfo

import "os"

func IsRoot() bool {
	return os.Getuid() == 0
}
