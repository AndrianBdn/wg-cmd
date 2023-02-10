package sysinfo

import (
	"fmt"
	"testing"
)

// Warning: this test goes to the Internet, actually we test that IP discovery ipDiscoveryServices still work
func TestDiscoverIP(t *testing.T) {
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
