package sysinfo

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestDiscoverIP(t *testing.T) {
	if os.Getenv("WG_TEST_HTTP") == "" {
		t.Skip("Skipping test that goes to the Internet")
	}

	for _, service := range ipDiscoveryServices {
		fmt.Printf("- Using %s to determine our IP\n", service)
		strIP, err := getMyIPWithService(service)
		if err == nil {
			fmt.Println("- Got response:", strIP)
		} else {
			t.Fatalf("Got error while using %s: %v", service, err)
		}
	}
}

func TestDiscoverIPFallback(t *testing.T) {
	saved := ipDiscoveryServices
	ipDiscoveryServices = []string{"http://127.0.0.1:1/"} // nothing listens there
	defer func() { ipDiscoveryServices = saved }()

	d := NewDiscoverIPStep()
	for i := 0; d.Result == "" && i < 10; i++ {
		d = DiscoverIP(d)
	}
	if !d.Fallback {
		t.Fatal("expected Fallback after all services fail")
	}
	if net.ParseIP(d.Result) == nil {
		t.Fatalf("fallback Result must be an IP address, got %q", d.Result)
	}
}

func TestDefaultRouteIP4(t *testing.T) {
	// sends no packets; only asks the kernel to resolve the default route
	ip, err := defaultRouteIP4()
	if err != nil {
		t.Skip("no default route:", err)
	}
	parsed := net.ParseIP(ip)
	if parsed == nil || parsed.To4() == nil {
		t.Fatalf("defaultRouteIP4 returned bad IP %q", ip)
	}
}
