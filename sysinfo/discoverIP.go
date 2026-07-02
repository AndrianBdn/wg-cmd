package sysinfo

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

var ipDiscoveryServices = []string{"https://ip4only.me/api/", "https://api.ipify.org/", "https://ifconfig.co/ip"}

type DiscoverStep struct {
	Result   string
	Fallback bool // Result is a local (default-route) address, not an externally verified IP
	step     int
	Service  string
	Log      string
}

func NewDiscoverIPStep() DiscoverStep {
	return DiscoverStep{
		step:    0,
		Service: ipDiscoveryServices[0],
	}
}

// DiscoverIP tries one service per call so the caller can render progress
// between attempts; feed the returned step back in until Result is set.
// When every service fails it falls back to the default-route local address
// (correct for a server with its public IP on the interface; wrong behind
// NAT) and sets Fallback so the UI can warn the user to verify the endpoint.
func DiscoverIP(d DiscoverStep) DiscoverStep {
	if d.Result != "" || d.step >= len(ipDiscoveryServices) {
		return d
	}

	strIP, err := getMyIPWithService(ipDiscoveryServices[d.step])
	if err == nil {
		d.Result = strIP
		return d
	}

	d.Log = err.Error()
	d.step++
	if d.step >= len(ipDiscoveryServices) {
		d.Service = ""
		d.Fallback = true
		if ip, err := defaultRouteIP4(); err == nil {
			d.Result = ip
		} else {
			d.Result = "127.0.0.1"
		}
		return d
	}
	d.Service = ipDiscoveryServices[d.step]
	return d
}

// defaultRouteIP4 returns the IPv4 source address the kernel picks for the
// default route. The UDP "connect" sends no packets — it only resolves
// routing, so it needs no Internet access, just a route.
func defaultRouteIP4() (string, error) {
	conn, err := net.Dial("udp4", "8.8.8.8:53")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || addr.IP.To4() == nil {
		return "", fmt.Errorf("cannot determine local IP4 address")
	}
	return addr.IP.String(), nil
}

func getMyIPWithService(serviceURL string) (string, error) {
	// as of January 2023 for the most people I know
	// can't use IP6 endpoint address for VPN
	// that's why we take special care do detect IPv4 address

	resp, err := ip4http().Get(serviceURL)
	if err != nil {
		return "", fmt.Errorf("getMyIPWithService http.get error: %w", err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == 200 {
		respBody, _ := io.ReadAll(resp.Body)
		sip := strings.TrimSpace(string(respBody))
		if strings.HasPrefix(sip, "IPv4,") {
			// https://ip4only.me/ API
			// IPv4,1.2.3.4,v1.1,,,See http://ip6.me/docs/ for api documentation
			parts := strings.Split(sip, ",")
			if len(parts) > 1 {
				sip = parts[1]
			}
		}

		ip := net.ParseIP(sip)
		if ip == nil {
			return "", fmt.Errorf("getMyIPWithService %s fail, bad IP %s", serviceURL, sip)
		}
		return ip.String(), nil
	}

	return "", fmt.Errorf("getMyIPWithService %s fail, bad response %s", serviceURL, resp.Status)
}

// addr is "host:port"
func resolveIPv4(addr string) (string, error) {
	url := strings.Split(addr, ":")
	if len(url) < 2 {
		return "", fmt.Errorf("bad addr")
	}
	ips, err := net.LookupIP(url[0])
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String() + ":" + url[1], nil
		}
	}
	return "", fmt.Errorf("no IP4")
}

func ip4http() *http.Client {
	// https://blog.bullgare.com/2021/02/force-ipv4-for-golang-http-client/
	// filled with defaults from go 1.19
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		ForceAttemptHTTP2: false,
		DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			ipv4, err := resolveIPv4(addr)
			if err != nil {
				return nil, err
			}

			return dialer.DialContext(ctx, network, ipv4)
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}
