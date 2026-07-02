package backend

import (
	"bytes"
	"strings"
	"testing"

	"gotest.tools/v3/golden"
)

func TestClient_GetPlainTextConfig(t *testing.T) {
	st, err := ReadState("testdata/wgc-wg3", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := st.Clients[11]
	cfg, err := c.GetPlainTextConfig(st.Server)
	if err != nil {
		t.Fatal(err)
	}
	golden.Assert(t, cfg, "wg3-dodo-client.conf.golden")
}

func TestNewClientPresharedKey(t *testing.T) {
	c1 := NewClient(2, "alice")
	c2 := NewClient(3, "bob")
	if c1.PresharedKey == "" {
		t.Fatal("NewClient must generate a per-peer PresharedKey")
	}
	if c1.PresharedKey == c2.PresharedKey {
		t.Fatal("per-peer PresharedKeys must be unique")
	}
}

func TestPeerPresharedKeyFallback(t *testing.T) {
	srv := &Server{PresharedKey: "server-wide"}

	legacy := &Client{} // peer file written before per-peer PSKs existed
	if k := srv.peerPresharedKey(legacy); k != "server-wide" {
		t.Fatalf("legacy peer must fall back to the server-wide key, got %q", k)
	}

	own := &Client{PresharedKey: "per-peer"}
	if k := srv.peerPresharedKey(own); k != "per-peer" {
		t.Fatalf("peer with own key must use it, got %q", k)
	}
}

func TestClientCreatedAtRoundtrip(t *testing.T) {
	dir := t.TempDir()

	c := NewClient(2, "alice")
	if c.CreatedAt.IsZero() {
		t.Fatal("NewClient must stamp CreatedAt")
	}
	if err := c.WriteOnce(dir); err != nil {
		t.Fatal(err)
	}

	rc, err := ReadClient(dir, c.GetFileName(), 2, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if !rc.CreatedAt.Equal(c.CreatedAt) {
		t.Fatalf("CreatedAt roundtrip mismatch: %v != %v", rc.CreatedAt, c.CreatedAt)
	}

	// legacy peer files have no CreatedAt and must decode to zero time
	st, err := ReadState("testdata/wgc-wg3", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Clients[11].CreatedAt.IsZero() {
		t.Fatal("legacy peer file must not have CreatedAt")
	}
}

func TestGeneratedConfigsUsePeerPSK(t *testing.T) {
	st, err := ReadState("testdata/wgc-wg3", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(12, "newpeer")

	clientCfg, err := c.GetPlainTextConfig(st.Server)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(clientCfg, "PresharedKey = "+c.PresharedKey) {
		t.Fatal("client config must contain the peer's own PresharedKey")
	}

	var buf bytes.Buffer
	if err := st.Server.generateServerPeerConfig(c, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "PresharedKey = "+c.PresharedKey) {
		t.Fatal("server peer block must contain the peer's own PresharedKey")
	}
}
