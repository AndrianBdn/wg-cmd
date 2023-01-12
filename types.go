package main

import (
	"encoding/hex"
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
	Interface        string
	ServerConfigPath string
	Address4         string
	Address6         string
	PresharedKey     string
	PrivateKey       string
	PublicKey        string
	ListenPort       uint16
	PostUp4          string
	PostDown4        string
	PostUp6          string
	PostDown6        string
	ClientRoute      string
	ClientDNS        string
	ServerHost       string
	addrInfo4        *addressInfo4
	addrInfo6        *addressInfo6
}

type addressInfo4 struct {
	server uint16
	prefix string
}

type addressInfo6 struct {
	server uint16
	prefix string
}

type client struct {
	ipNum      int
	name       string
	fileName   string
	PublicKey  string
	PrivateKey string
}

type state struct {
	server  *server
	clients map[int]*client
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
			return nil, fmt.Errorf("address4 only /24 CIDR supported")
		}
		prefix := strconv.Itoa(int(ip4[0])) + "." +
			strconv.Itoa(int(ip4[1])) + "." +
			strconv.Itoa(int(ip4[2])) + "."
		s.addrInfo4 = &addressInfo4{prefix: prefix, server: uint16(ip4[3])}
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

		ip6parts := strings.SplitAfter(s.Address6, "::")
		if len(ip6parts) != 2 {
			return nil, fmt.Errorf("address6 cannot find :: in IP6 addr")
		}
		s.addrInfo6 = &addressInfo6{prefix: ip6parts[0], server: uint16(ip[15])}
	}

	return &s, nil
}

func readClient(fileName string, ipNum int, name string) (*client, error) {
	var c client
	_, err := toml.DecodeFile(fileName, &c)
	if err != nil {
		return nil, fmt.Errorf("readClient error %w", err)
	}
	c.fileName = fileName
	c.ipNum = ipNum
	c.name = name
	return &c, nil
}

func newClient(ip int, name string) *client {
	fileName := fmt.Sprintf("%03d-%s.toml", ip, name)
	c := client{ipNum: ip, name: name, fileName: fileName}
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic("Can't generate wireguard keypair, err: " + err.Error())
	}
	c.PrivateKey = key.String()
	c.PublicKey = key.PublicKey().String()
	return &c
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

	return &s
}

func writeConfigHeader(f *os.File) {
	_, _ = f.WriteString(strings.TrimLeft(dedent.Dedent(`
				# Warning: this is not a Wireguard config
				# This file uses TOML (toml.io) syntax, instead of Wireguard 
				# wg-dir-conf tool builds Wireguard config using files in this directories
				# >> You are welcome to edit this file, it wont be overwritten.
				`), "\n"))
}

func (s *server) save() error {
	f, err := os.OpenFile(serverFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	writeConfigHeader(f)
	if err != nil {
		return fmt.Errorf("server.save, can't create %s file %w", serverFileName, err)
	}

	if err := toml.NewEncoder(f).Encode(s); err != nil {
		return fmt.Errorf("server.save, error TOML encoding server struct %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("server.save, can't close %s file %w", serverFileName, err)
	}

	return nil
}

func (c *client) writeOnce() error {
	f, err := os.OpenFile(c.fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	writeConfigHeader(f)
	if err != nil {
		return fmt.Errorf("client.writeOnce, can't create %s file %w", c.fileName, err)
	}

	if err := toml.NewEncoder(f).Encode(c); err != nil {
		return fmt.Errorf("client.writeOnce, error TOML encoding server struct %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("client.writeOnce, can't close %s file %w", c.fileName, err)
	}

	return nil
}

func (c *client) allowedIps(srv *server) (string, error) {
	result := ""
	if srv.Address4 != "" {
		if int(srv.addrInfo4.server) == c.ipNum {
			return "", fmt.Errorf("client x.x.x.%d IP (filename) is set to the same as server %s", c.ipNum, srv.Address4)
		}
		result += srv.addrInfo4.prefix + strconv.Itoa(c.ipNum) + "/32"
	}
	if srv.Address6 != "" {
		if int(srv.addrInfo6.server) == c.ipNum {
			return "", fmt.Errorf("client x:x:x::%d IP (filename) is set to the same as server %s", c.ipNum, srv.Address6)
		}
		if result != "" {
			result += ", "
		}
		var lb [1]byte
		lb[0] = byte(c.ipNum)
		result += srv.addrInfo6.prefix + hex.EncodeToString(lb[:]) + "/64"
	}
	return result, nil
}
