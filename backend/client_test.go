package backend

import (
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
