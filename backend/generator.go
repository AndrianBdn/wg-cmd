package backend

import (
	"fmt"
	"io"
	"strconv"
)

func (s *State) generateWireguardConfig(w io.Writer) error {
	err := s.Server.WriteInterfaceBlock(w, true)
	if err != nil {
		return fmt.Errorf("WriteInterfaceBlock error %w", err)
	}

	return s.IterateClients(func(c *Client) error {
		err = s.Server.generateServerPeerConfig(c, w)
		if err != nil {
			return fmt.Errorf("generateServerPeerConfig error %w", err)
		}
		return nil
	})
}

func (s *Server) generateServerPeerConfig(client *Client, w io.Writer) error {
	_, err := fmt.Fprintf(w, "\n# peer %s\n", client.name)
	if err != nil {
		return fmt.Errorf("generateServerConfig error %w", err)
	}
	_, _ = fmt.Fprintln(w, "[Peer]")
	if s.PresharedKey != "" {
		_, _ = fmt.Fprintln(w, "PresharedKey =", s.PresharedKey)
	}
	_, _ = fmt.Fprintln(w, "PublicKey =", client.PublicKey)
	_, _ = fmt.Fprintln(w, "AllowedIPs =", client.AllowedIps(s))

	return nil
}

func (c *Client) generateClientConfig(server *Server, w io.Writer) error {
	_, err := fmt.Fprintf(w, "[Interface]\n")
	if err != nil {
		return fmt.Errorf("GenerateClientConfig error %w", err)
	}
	_, _ = fmt.Fprintln(w, "PrivateKey =", c.PrivateKey)
	_, _ = fmt.Fprintln(w, "Address =", c.AllowedIps(server))

	if server.ClientDNS != "" {
		_, _ = fmt.Fprintln(w, "DNS =", server.ClientDNS)
	}

	_, _ = fmt.Fprintln(w, "\n[Peer]")
	if server.PresharedKey != "" {
		_, _ = fmt.Fprintln(w, "PresharedKey =", server.PresharedKey)
	}
	_, _ = fmt.Fprintln(w, "PublicKey =", server.PublicKey)
	cr := c.ClientRoute
	if cr == "" {
		cr = server.ClientRoute
	}

	_, _ = fmt.Fprintln(w, "AllowedIPs =", cr)
	_, _ = fmt.Fprintln(w, "Endpoint =", server.ClientServerEndpoint+":"+strconv.Itoa(int(server.ListenPort)))

	if server.ClientPersistentKeepalive != 0 {
		_, _ = fmt.Fprintln(w, "PersistentKeepalive =", server.ClientPersistentKeepalive)
	}

	return nil
}
