//go:build linux

package sysinfo

import (
	"os"
	"strconv"
	"strings"
)

type procNetRouteLine struct {
	iface       string
	mask        uint32
	destination uint32
}

type procNetIPV6RouteLine struct {
	iface    string
	metric   int32
	destZero bool
}

func NetworkInterfaceExists(iface string) bool {
	ifls4, err := readProcNetRoute()
	if err == nil {
		for _, ifl := range ifls4 {
			if ifl.iface == iface {
				return true
			}
		}
	}

	ifls6, err := readProcNetIPV6Route()
	if err == nil {
		for _, ifl := range ifls6 {
			if ifl.iface == iface {
				return true
			}
		}
	}

	return false
}

func DefaultIP4Interface() string {
	ifls, err := readProcNetRoute()
	if err != nil {
		return defaultUnknownIface
	}

	for _, ifl := range ifls {
		if ifl.mask == 0 && ifl.destination == 0 {
			return ifl.iface
		}
	}

	return defaultUnknownIface
}

func HasIP6() bool {
	return ip6interface() != ""
}

func DefaultIP6Interface() string {
	r := ip6interface()
	if r == "" {
		return defaultUnknownIface
	}
	return r
}

func ip6interface() string {
	ifls, err := readProcNetIPV6Route()
	if err != nil {
		return ""
	}

	for _, ifl := range ifls {
		if ifl.metric > 0 && ifl.destZero {
			return ifl.iface
		}
	}
	return ""
}

func readProcNetRoute() ([]procNetRouteLine, error) {
	r, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(r), "\n")
	ret := make([]procNetRouteLine, 0, 16)
	for _, line := range lines {
		if strings.HasPrefix(line, "Iface") {
			continue // header
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 10 {
			continue
		}
		// here we don't care about byte order
		mask, err := strconv.ParseUint(fields[7], 16, 32)
		if err != nil {
			return nil, err
		}
		dest, err := strconv.ParseUint(fields[1], 16, 32)
		if err != nil {
			return nil, err
		}
		sl := procNetRouteLine{fields[0], uint32(mask), uint32(dest)}
		ret = append(ret, sl)
	}
	return ret, nil
}

func readProcNetIPV6Route() ([]procNetIPV6RouteLine, error) {
	r, err := os.ReadFile("/proc/net/ipv6_route")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(r), "\n")
	ret := make([]procNetIPV6RouteLine, 0, 16)
	for _, line := range lines {
		fields := strings.Split(line, " ")
		if len(fields) < 10 {
			continue
		}

		destZero := fields[0] == "00000000000000000000000000000000"
		metric, err := strconv.ParseUint(fields[5], 16, 32)
		if err != nil {
			return nil, err
		}

		// here we don't care about byte order
		sl := procNetIPV6RouteLine{iface: fields[len(fields)-1], destZero: destZero, metric: int32(metric)}
		ret = append(ret, sl)
	}
	return ret, nil
}
