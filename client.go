package main

import (
	"encoding/hex"
	"fmt"
	"github.com/BurntSushi/toml"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"os"
	"strconv"
)

type client struct {
	ipNum      int
	name       string
	fileName   string
	PublicKey  string
	PrivateKey string
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

func (c *client) allowedIps(srv *server) string {
	result := ""
	if srv.Address4 != "" {
		result += srv.addrInfo4.prefix + strconv.Itoa(c.ipNum) + "/32"
	}
	if srv.Address6 != "" {
		if result != "" {
			result += ", "
		}
		var lb [1]byte
		lb[0] = byte(c.ipNum)
		result += srv.addrInfo6.prefix + hex.EncodeToString(lb[:]) + "/128"
	}
	return result
}
