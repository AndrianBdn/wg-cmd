//go:build !linux

package sysinfo

func DefaultIP4Interface() string {
	return defaultUnknownIface
}

func HasIP6() bool {
	return true
}

func DefaultIP6Interface() string {
	return defaultUnknownIface
}

func NetworkInterfaceExists(iface string) bool {
	return false
}
