package backend

import (
	"os"
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

// TestState_EmailPeerNameRoundTrip makes sure an email-style peer name (dots
// and an @) survives being written to a NNNN-<name>.toml file and read back
// from disk, and that path-unsafe names are still rejected.
func TestState_EmailPeerNameRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS("testdata/wgc-wg3")); err != nil {
		t.Fatal(err)
	}

	s, err := ReadState(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	const email = "john.doe@example.com"
	if err := s.AddPeer(email); err != nil {
		t.Fatalf("AddPeer(%q): %v", email, err)
	}

	// path separators must never be accepted (would escape the interface dir)
	if err := s.AddPeer("john/../../etc"); err == nil {
		t.Fatal("expected error for name with slash, got nil")
	}

	// consecutive dots are rejected (defense-in-depth; no valid email has them)
	if err := s.AddPeer("john..doe@example.com"); err == nil {
		t.Fatal("expected error for name with consecutive dots, got nil")
	}

	// email:device convention (colon) is allowed and must round-trip too
	const device = "bob@corp.com:laptop"
	if err := s.AddPeer(device); err != nil {
		t.Fatalf("AddPeer(%q): %v", device, err)
	}

	// re-read from disk: the file names must parse back to the exact names
	s2, err := ReadState(dir, nil)
	if err != nil {
		t.Fatalf("re-read state: %v", err)
	}

	for name, wantFile := range map[string]string{
		email:  "0012-" + email + ".toml",
		device: "0013-" + device + ".toml",
	} {
		var found *Client
		for _, c := range s2.Clients {
			if c.GetName() == name {
				found = c
				break
			}
		}
		if found == nil {
			t.Fatalf("peer %q not found after re-reading state from disk", name)
		}
		if found.GetFileName() != wantFile {
			t.Fatalf("expected file %q, got %q", wantFile, found.GetFileName())
		}
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
