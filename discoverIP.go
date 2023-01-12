package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func discoverIP() string {
	services := []string{"http://whatismyip.akamai.com/", "https://api.ipify.org/", "https://ifconfig.co/ip"}
	for _, service := range services {
		fmt.Printf("- Using %s to determine our IP\n", service)
		strIP, err := getMyIPWithService(service)
		if err == nil {
			fmt.Println("- Got response:", strIP)
			return strIP
		}
	}
	fmt.Println("- All services failed, returning 127.0.0.1")
	return "127.0.0.1"
}

func getMyIPWithService(serviceURL string) (string, error) {
	resp, err := http.Get(serviceURL)
	if err != nil {
		return "", fmt.Errorf("getMyIPWithService http.get error: %w", err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == 200 {
		respBody, _ := io.ReadAll(resp.Body)
		sip := strings.TrimSpace(string(respBody))
		ip := net.ParseIP(sip)
		if ip == nil {
			return "", fmt.Errorf("getMyIPWithService %s fail, bad IP %s", serviceURL, sip)
		}
		return ip.String(), nil
	}

	return "", fmt.Errorf("getMyIPWithService %s fail, bad response %s", serviceURL, resp.Status)
}
