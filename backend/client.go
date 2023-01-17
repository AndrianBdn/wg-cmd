package backend

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/BurntSushi/toml"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"os"
	"strconv"
)

type Client struct {
	ipNum      int
	name       string
	fileName   string
	PublicKey  string
	PrivateKey string
}

func ReadClient(fileName string, ipNum int, name string) (*Client, error) {
	var c Client
	_, err := toml.DecodeFile(fileName, &c)
	if err != nil {
		return nil, fmt.Errorf("ReadClient error %w", err)
	}
	c.fileName = fileName
	c.ipNum = ipNum
	c.name = name
	return &c, nil
}

func NewClient(ip int, name string) *Client {
	fileName := fmt.Sprintf("%03d-%s.toml", ip, name)
	c := Client{ipNum: ip, name: name, fileName: fileName}
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic("Can't generate wireguard keypair, err: " + err.Error())
	}
	c.PrivateKey = key.String()
	c.PublicKey = key.PublicKey().String()
	return &c
}

func (c *Client) WriteOnce() error {
	f, err := os.OpenFile(c.fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	writeConfigHeader(f)
	if err != nil {
		return fmt.Errorf("Client.WriteOnce, can't create %s file %w", c.fileName, err)
	}

	if err := toml.NewEncoder(f).Encode(c); err != nil {
		return fmt.Errorf("Client.WriteOnce, error TOML encoding Server struct %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("Client.WriteOnce, can't close %s file %w", c.fileName, err)
	}

	return nil
}

func (c *Client) AllowedIps(srv *Server) string {
	result := ""
	if srv.Address4 != "" {
		cAdd := c.ipNum / 256
		d := c.ipNum % 256
		result += srv.addrInfo4.prefix + strconv.Itoa(int(srv.addrInfo4.c)+cAdd) + "." + strconv.Itoa(d) + "/32"
	}
	if srv.Address6 != "" {
		if result != "" {
			result += ", "
		}
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(c.ipNum))
		result += srv.addrInfo6.prefix + hex.EncodeToString(b) + "/128"
	}
	return result
}

func (c *Client) GetName() string {
	return c.name
}

func (c *Client) GetFileName() string {
	return c.fileName
}

func (c *Client) GetIPNumber() int {
	return c.ipNum
}
