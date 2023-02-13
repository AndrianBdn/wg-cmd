package sysinfo

import (
	"fmt"
	"os"
	"testing"
)

func TestDiscoverIP(t *testing.T) {
	if os.Getenv("WG_TEST_HTTP") == "" {
		t.Skip("Skipping test that goes to the Internet")
	}

	for _, service := range ipDiscoveryServices {
		if service == "" {
			continue
		}
		fmt.Printf("- Using %s to determine our IP\n", service)
		strIP, err := getMyIPWithService(service)
		if err == nil {
			fmt.Println("- Got response:", strIP)
		} else {
			t.Fatalf("Got error while using %s: %v", service, err)
		}
	}
}
