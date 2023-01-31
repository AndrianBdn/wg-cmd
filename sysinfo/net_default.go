//go:build !linux

package sysinfo

func DefaultIP4Interface() string {
	return defaultUnknownIface
}

func HasIP6() bool {
	return false
}

func DefaultIP6Interface() string {
	return defaultUnknownIface
}

func NetworkInterfaceExists(iface string) bool {
	return false
}
