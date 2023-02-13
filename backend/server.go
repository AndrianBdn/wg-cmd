package backend

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/andrianbdn/wg-cmd/sysinfo"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const ServerFileName = `0001-server.toml`

type ServerBlueprint struct {
	InterfaceName string
	Endpoint      string
	Port          uint16
	Nat4          bool
	Nat6          bool
	Net4          string
	Net6          string
	DNS           string
}

type Server struct {
	Interface                 string
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
	ClientServerEndpoint      string
	ClientPersistentKeepalive int
	addrInfo4                 *addressInfo4
	addrInfo6                 *addressInfo6
}

type addressInfo4 struct {
	prefix string
	c      uint8
}

type addressInfo6 struct {
	prefix string
}

func NewServerWithBlueprint(b ServerBlueprint) *Server {
	s := Server{}
	s.Address4 = netToServerIP4(b.Net4)
	if b.Net6 != "" {
		s.Address6 = netToServerIP6(b.Net6)
	}
	s.Interface = b.InterfaceName
	s.ListenPort = b.Port
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

	if b.Nat4 {
		s.PostUp4 = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; "+
			"iptables -t nat -A POSTROUTING -o %s -j MASQUERADE", b.InterfaceName, sysinfo.DefaultIP4Interface())

		s.PostDown4 = fmt.Sprintf("iptables -D FORWARD -i %s -j ACCEPT; "+
			"iptables -t nat -D POSTROUTING -o %s -j MASQUERADE", b.InterfaceName, sysinfo.DefaultIP4Interface())
	}

	if b.Nat6 {
		s.PostUp6 = fmt.Sprintf("ip6tables -A FORWARD -i %s -j ACCEPT; "+
			"ip6tables -t nat -A POSTROUTING -o %s -j MASQUERADE", b.InterfaceName, sysinfo.DefaultIP6Interface())
		s.PostDown6 = fmt.Sprintf("ip6tables -D FORWARD -i %s -j ACCEPT; "+
			"ip6tables -t nat -D POSTROUTING -o %s -j MASQUERADE", b.InterfaceName, sysinfo.DefaultIP6Interface())
	}

	if b.Nat4 {
		s.ClientRoute = "0.0.0.0/0"
	} else {
		s.ClientRoute = b.Net4
	}

	if b.Net6 != "" {
		s.ClientRoute += ", "
		if b.Nat6 {
			s.ClientRoute += "::/0"
		} else {
			s.ClientRoute += b.Net6
		}
	}

	s.ClientDNS = b.DNS
	s.ClientServerEndpoint = b.Endpoint
	s.ClientPersistentKeepalive = 42
	return &s
}

func NewServer(iface string, serverHost string) *Server {
	b := ServerBlueprint{}
	b.Net4 = RandomIP4("")
	b.Net6 = RandomIP6()
	b.InterfaceName = iface
	b.Endpoint = serverHost
	b.Port = 51820
	b.Nat4 = true
	b.Nat6 = true
	b.DNS = "1.1.1.1"
	return NewServerWithBlueprint(b)
}

func ReadServer() (*Server, error) {
	var s Server
	_, err := toml.DecodeFile(ServerFileName, &s)
	if err != nil {
		return nil, fmt.Errorf("ReadServer error %w", err)
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
		if ones != 20 {
			return nil, fmt.Errorf("address4 supports only /24 network")
		}
		if ip4[3] != 1 {
			return nil, fmt.Errorf("address4 Server IP must start with 1")
		}
		prefix := strconv.Itoa(int(ip4[0])) + "." +
			strconv.Itoa(int(ip4[1])) + "."
		s.addrInfo4 = &addressInfo4{prefix: prefix, c: ip4[2]}
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
			return nil, fmt.Errorf("address6 Server IP must start with 1")
		}
		ip6parts := strings.SplitAfter(s.Address6, "::")
		if len(ip6parts) != 2 {
			return nil, fmt.Errorf("address6 cannot find :: in IP6 addr")
		}
		s.addrInfo6 = &addressInfo6{prefix: ip6parts[0]}
	}

	return &s, nil
}

func (s *Server) WriteOnce() error {
	f, err := os.OpenFile(ServerFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	writeConfigHeader(f)
	if err != nil {
		return fmt.Errorf("Server.WriteOnce, can't create %s file %w", ServerFileName, err)
	}

	if err := toml.NewEncoder(f).Encode(s); err != nil {
		return fmt.Errorf("Server.WriteOnce, error TOML encoding Server struct %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("Server.WriteOnce, can't close %s file %w", ServerFileName, err)
	}

	return nil
}

func (s *Server) WriteInterfaceBlock(w io.Writer, writeComment bool) error {
	PostUp := ""
	PostDown := ""

	comment := ""
	if writeComment {
		wd, err := os.Getwd()
		if err != nil {
			panic("can't os.Getwd " + err.Error())
		}
		comment := "# This file is generated by WG Commander from directory " + wd
		comment += "\n# It it likely to be overwritten.\n\n"
	}

	_, err := fmt.Fprintf(w, "%s[Interface]\n", comment)
	if err != nil {
		return fmt.Errorf("generateServerConfig error %w", err)
	}
	if s.Address4 != "" {
		_, _ = fmt.Fprintln(w, "Address =", s.Address4)
		PostUp = strings.TrimRight(s.PostUp4, " ;")
		PostDown = strings.TrimRight(s.PostDown4, "; ")
	}
	if s.Address6 != "" {
		_, _ = fmt.Fprintln(w, "Address =", s.Address6)
		PostUp = concatIfNotEmpty(PostUp, "; ")
		PostDown = concatIfNotEmpty(PostDown, "; ")

		PostUp = PostUp + s.PostUp6
		PostDown = PostDown + s.PostDown6
	}
	_, _ = fmt.Fprintln(w, "PostUp =", PostUp)
	_, _ = fmt.Fprintln(w, "PostDown =", PostDown)
	_, _ = fmt.Fprintln(w, "ListenPort =", s.ListenPort)
	_, _ = fmt.Fprintln(w, "PrivateKey =", s.PrivateKey)
	return nil
}

func (s *Server) GetInterfaceString() string {
	var b strings.Builder
	err := s.WriteInterfaceBlock(&b, false)
	if err != nil {
		panic("GetInterfaceString error " + err.Error())
	}
	return b.String()
}
