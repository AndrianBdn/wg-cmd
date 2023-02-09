package backend

import (
	"fmt"
	"io"
	"strconv"
)

func generateWireguardConfig(state *State, w io.Writer) error {
	err := state.Server.WriteInterfaceBlock(w, true)
	if err != nil {
		return fmt.Errorf("WriteInterfaceBlock error %w", err)
	}

	for _, client := range state.Clients {
		err = generateServerPeerConfig(state.Server, client, w)
		if err != nil {
			return fmt.Errorf("generateServerConfig error %w", err)
		}
	}

	return nil
}

func generateServerPeerConfig(srv *Server, client *Client, w io.Writer) error {
	_, err := fmt.Fprintf(w, "\n# peer %s\n", client.name)
	if err != nil {
		return fmt.Errorf("generateServerConfig error %w", err)
	}
	_, _ = fmt.Fprintln(w, "[Peer]")
	if srv.PresharedKey != "" {
		_, _ = fmt.Fprintln(w, "PresharedKey =", srv.PresharedKey)
	}
	_, _ = fmt.Fprintln(w, "PublicKey =", client.PublicKey)
	_, _ = fmt.Fprintln(w, "AllowedIPs =", client.AllowedIps(srv))

	return nil
}

func GenerateClientConfig(server *Server, client *Client, w io.Writer) error {
	_, err := fmt.Fprintf(w, "[Interface]\n")
	if err != nil {
		return fmt.Errorf("GenerateClientConfig error %w", err)
	}
	_, _ = fmt.Fprintln(w, "PrivateKey =", client.PrivateKey)
	_, _ = fmt.Fprintln(w, "Address =", client.AllowedIps(server))

	if server.ClientDNS != "" {
		_, _ = fmt.Fprintln(w, "DNS =", server.ClientDNS)
	}

	_, _ = fmt.Fprintln(w, "\n[Peer]")
	if server.PresharedKey != "" {
		_, _ = fmt.Fprintln(w, "PresharedKey =", server.PresharedKey)
	}
	_, _ = fmt.Fprintln(w, "PublicKey =", server.PublicKey)
	_, _ = fmt.Fprintln(w, "AllowedIPs =", server.ClientRoute)
	_, _ = fmt.Fprintln(w, "Endpoint =", server.ClientServerEndpoint+":"+strconv.Itoa(int(server.ListenPort)))

	if server.ClientPersistentKeepalive != 0 {
		_, _ = fmt.Fprintln(w, "PersistentKeepalive =", server.ClientPersistentKeepalive)
	}

	return nil
}
