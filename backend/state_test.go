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
