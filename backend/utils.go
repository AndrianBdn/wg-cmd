package backend

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/lithammer/dedent"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
)

func writeConfigHeader(f *os.File) {
	_, _ = f.WriteString(strings.TrimLeft(dedent.Dedent(`
				# Warning: this is not a Wireguard config
				# This file uses TOML (toml.io) syntax, instead of Wireguard 
				# wg-dir-conf tool builds Wireguard config using files in this directory
				# >> You are welcome to edit this file, it wont be overwritten.
				`), "\n"))
}

func RandomIP4() string {
	return "10." + strconv.Itoa(mrand.Intn(256)) + "." + strconv.Itoa(mrand.Intn(16)*16) + ".0/20"
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
