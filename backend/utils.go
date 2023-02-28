package backend

import (
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/lithammer/dedent"
)

func writeConfigHeader(f *os.File) {
	_, _ = f.WriteString(strings.TrimLeft(dedent.Dedent(`
				# Warning: this is not a WireGuard config
				# This file uses TOML (toml.io) syntax, instead of WireGuard 
				# WG Commander builds WireGuard config using files in this directory
				# >> You are welcome to edit this file, it wont be overwritten.
				`), "\n"))
}

func RandomIP4(A string) string {
	ab := ""

	if A == "100" {
		ab = "100." + strconv.Itoa(64+mrand.Intn(64))
	} else if A == "172" {
		ab = "172." + strconv.Itoa(16+mrand.Intn(16))
	} else if A == "192" {
		ab = "192.168"
	} else {
		ab = "10." + strconv.Itoa(mrand.Intn(256))
	}

	return ab + "." + strconv.Itoa(mrand.Intn(16)*16) + ".0/20"
}

func RandomIP6() string {
	// I don't think that RFC 4193 about SHA1 and machine id is sane,
	// so let's just read cryptographically random bytes
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic("failed to read 5 random bytes for IP6" + err.Error())
	}
	bhex := hex.EncodeToString(b)
	return "fd" + bhex[0:2] + ":" + bhex[2:6] + ":" + bhex[6:10] + "::0/64"
}

func netToServerIP4(ip4 string) string {
	return strings.ReplaceAll(ip4, ".0/", ".1/")
}

func netToServerIP6(ip4 string) string {
	return strings.ReplaceAll(ip4, "::0/", "::1/")
}

func concatIfNotEmpty(str string, add string) string {
	if str != "" {
		return str + add
	}
	return str
}

func dnsString(dns []string, ipv6 bool) string {
	var dnsStr string
	for _, d := range dns {
		if strings.Contains(d, ":") && !ipv6 {
			continue
		}
		dnsStr = concatIfNotEmpty(dnsStr, ", ")
		dnsStr += d
	}
	return dnsStr
}
