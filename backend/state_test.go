package backend

import (
	"strings"
	"testing"

	"gotest.tools/v3/golden"
)

func TestStateRead(t *testing.T) {
	s, err := ReadState("testdata/wgc-wg3", nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(s.Clients) != 10 {
		t.Fatal("wrong number of clients")
	}
}

func TestState_AddPeer(t *testing.T) {
	s, err := ReadState("testdata/wgc-wg3", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddPeer("test")
	if err != nil {
		t.Fatal(err)
	}

	if len(s.Clients) != 11 {
		t.Fatal("wrong number of clients")
	}

	err = s.DeletePeer(12)
	if err != nil {
		t.Fatal(err)
	}
}

func TestState_RenamePeer(t *testing.T) {
	s, err := ReadState("testdata/wgc-wg3", nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.AddPeer("test"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.DeletePeer(12); err != nil {
			t.Fatal(err)
		}
	}()

	if err := s.RenamePeer(12, "renamed"); err != nil {
		t.Fatal(err)
	}

	if s.Clients[12].GetName() != "renamed" {
		t.Fatalf("expected name renamed, got %s", s.Clients[12].GetName())
	}
	if s.Clients[12].GetFileName() != "0012-renamed.toml" {
		t.Fatalf("expected file 0012-renamed.toml, got %s", s.Clients[12].GetFileName())
	}

	// renaming to a name already in use must fail
	if err := s.RenamePeer(12, "alice"); err == nil {
		t.Fatal("expected error renaming to existing name, got nil")
	}

	// invalid names must be rejected
	if err := s.RenamePeer(12, "bad/name"); err == nil {
		t.Fatal("expected error renaming to invalid name, got nil")
	}
}

func TestState_GenerateWireguardFile(t *testing.T) {
	s, err := ReadState("testdata/wgc-wg3", nil)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(strings.Builder)
	err = s.generateWireguardConfig(buf)
	if err != nil {
		t.Fatal(err)
	}

	golden.Assert(t, buf.String(), "wg3.conf.golden")
}
