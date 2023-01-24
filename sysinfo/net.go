package sysinfo

func DefaultIP4Interface() string {
	return "en0"
}

func HasIP6() bool {
	return false
}

func DefaultIP6Interface() string {
	return "en0"
}
