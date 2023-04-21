package backend

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Client struct {
	ipNum          int
	name           string
	fileName       string
	PublicKey      string
	PrivateKey     string
	ClientRoute    string
	AddServerRoute string
	MTU            int
	DNS            string
}

func ReadClient(dir string, fileName string, ipNum int, name string) (*Client, error) {
	var c Client
	_, err := toml.DecodeFile(filepath.Join(dir, fileName), &c)
	if err != nil {
		return nil, fmt.Errorf("ReadClient %w", err)
	}
	c.fileName = fileName
	c.ipNum = ipNum
	c.name = name
	return &c, nil
}

func NewClient(ip int, name string) *Client {
	fileName := fmt.Sprintf("%04d-%s.toml", ip, name)
	c := Client{ipNum: ip, name: name, fileName: fileName}
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic("Can't generate wireguard keypair, err: " + err.Error())
	}
	c.PrivateKey = key.String()
	c.PublicKey = key.PublicKey().String()
	return &c
}

func (c *Client) WriteOnce(dir string) error {
	f, err := os.OpenFile(filepath.Join(dir, c.fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
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
		result += c.GetIP4(srv)
	}
	if srv.Address6 != "" {
		result = concatIfNotEmpty(result, ", ")
		result += c.GetIP6(srv)
	}
	return result
}

func (c *Client) GetIP4(srv *Server) string {
	cAdd := c.ipNum / 256
	d := c.ipNum % 256
	return srv.addrInfo4.prefix + strconv.Itoa(int(srv.addrInfo4.c)+cAdd) + "." + strconv.Itoa(d) + "/32"
}

func (c *Client) GetIP6(srv *Server) string {
	if srv.addrInfo6 == nil {
		return ""
	}
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(c.ipNum))
	return srv.addrInfo6.prefix + hex.EncodeToString(b) + "/128"
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

func (c *Client) GetIPNumberString() string {
	return fmt.Sprintf("%04d", c.ipNum)
}

func (c *Client) GetPlainTextConfig(srv *Server) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := c.generateClientConfig(srv, buf)
	if err != nil {
		return "", fmt.Errorf("generate client config %w", err)
	}

	return buf.String(), nil
}
