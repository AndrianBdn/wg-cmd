//go:build !linux

package sysinfo

import "errors"

func HasSysctl() bool {
	return false
}

func NatEnabledIPv4() bool {
	return false
}

func NatEnabledIPv6() bool {
	return false
}

func EnableNat(ip4, ip6 bool) error {
	return errors.New("not implemented")
}

func HasSystemd() bool {
	return false
}

func CreateSystemdStuff(iface, wgdir string) error {
	return errors.New("not implemented")
}
