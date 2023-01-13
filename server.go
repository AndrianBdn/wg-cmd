package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/lithammer/dedent"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"os"
	"strconv"
	"strings"
)

const serverFileName = "001-server.toml"

type server struct {
	Interface                 string
	ServerConfigPath          string
	Address4                  string
	Address6                  string
	PresharedKey              string
	PrivateKey                string
	PublicKey                 string
	ListenPort                uint16
	PostUp4                   string
	PostDown4                 string
	PostUp6                   string
	PostDown6                 string
	ClientRoute               string
	ClientDNS                 string
	ServerHost                string
	ClientPersistentKeepalive int
	addrInfo4                 *addressInfo4
	addrInfo6                 *addressInfo6
}

type addressInfo4 struct {
	prefix string
}

type addressInfo6 struct {
	prefix string
}

func newServer(iface string, addr4 string, addr6 string, serverHost string) *server {
	s := server{}
	s.Address4 = addr4
	s.Address6 = addr6
	s.Interface = iface
	s.ServerConfigPath = "/etc/wireguard/" + iface + ".conf"
	s.ListenPort = 51820
	key, err := wgtypes.GenerateKey()
	if err != nil {
		panic("Can't generate wireguard pre-shared key, err: " + err.Error())
	}
	s.PresharedKey = key.String()

	key, err = wgtypes.GeneratePrivateKey()
	if err != nil {
		panic("Can't generate wireguard keypair, err: " + err.Error())
	}
	s.PrivateKey = key.String()
	s.PublicKey = key.PublicKey().String()
	s.PostUp4 = "iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE"
	s.PostUp6 = "ip6tables -A FORWARD -i wg0 -j ACCEPT; ip6tables -t nat -A POSTROUTING -o eth0 -j MASQUERADE"
	s.PostDown4 = "iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE"
	s.PostDown6 = "ip6tables -D FORWARD -i wg0 -j ACCEPT; ip6tables -t nat -D POSTROUTING -o eth0 -j MASQUERADE"
	s.ClientRoute = "0.0.0.0/0, ::/0"
	s.ClientDNS = "1.1.1.1"
	s.ServerHost = serverHost
	s.ClientPersistentKeepalive = 42
	return &s
}

func readServer() (*server, error) {
	var s server
	_, err := toml.DecodeFile(serverFileName, &s)
	if err != nil {
		return nil, fmt.Errorf("readServer error %w", err)
	}

	if s.Address4 != "" {
		ip, ipNet, err := net.ParseCIDR(s.Address4)
		if err != nil {
			return nil, fmt.Errorf("address4 parser error %w", err)
		}
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, fmt.Errorf("address4 must contain IP4 address")
		}
		ones, _ := ipNet.Mask.Size()
		if ones != 24 {
			return nil, fmt.Errorf("address4 supports only /24 network")
		}
		if ip4[3] != 1 {
			return nil, fmt.Errorf("address4 server IP must start with 1")
		}
		prefix := strconv.Itoa(int(ip4[0])) + "." +
			strconv.Itoa(int(ip4[1])) + "." +
			strconv.Itoa(int(ip4[2])) + "."
		s.addrInfo4 = &addressInfo4{prefix: prefix}
	}

	if s.Address6 != "" {
		ip, ipNet, err := net.ParseCIDR(s.Address6)
		if err != nil {
			return nil, fmt.Errorf("address6 parser error %w", err)
		}
		ip4 := ip.To4()
		if ip4 != nil {
			return nil, fmt.Errorf("address6 must contain IP6 address")
		}
		ones, _ := ipNet.Mask.Size()
		if ones != 64 {
			return nil, fmt.Errorf("address6 only /64 CIDR supported")
		}
		if ip[15] != 1 {
			return nil, fmt.Errorf("address6 server IP must start with 1")
		}
		ip6parts := strings.SplitAfter(s.Address6, "::")
		if len(ip6parts) != 2 {
			return nil, fmt.Errorf("address6 cannot find :: in IP6 addr")
		}
		s.addrInfo6 = &addressInfo6{prefix: ip6parts[0]}
	}

	return &s, nil
}

func (s *server) writeOnce() error {
	f, err := os.OpenFile(serverFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	writeConfigHeader(f)
	if err != nil {
		return fmt.Errorf("server.writeOnce, can't create %s file %w", serverFileName, err)
	}

	if err := toml.NewEncoder(f).Encode(s); err != nil {
		return fmt.Errorf("server.writeOnce, error TOML encoding server struct %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("server.writeOnce, can't close %s file %w", serverFileName, err)
	}

	return nil
}

func writeConfigHeader(f *os.File) {
	_, _ = f.WriteString(strings.TrimLeft(dedent.Dedent(`
				# Warning: this is not a Wireguard config
				# This file uses TOML (toml.io) syntax, instead of Wireguard 
				# wg-dir-conf tool builds Wireguard config using files in this directory
				# >> You are welcome to edit this file, it wont be overwritten.
				`), "\n"))
}
