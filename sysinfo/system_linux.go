//go:build linux

package sysinfo

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func HasSystemd() bool {
	systemdDir := "/etc/systemd/system"
	if _, err := os.Stat(systemdDir); err == nil {
		return true
	}
	return false
}

func HasSysctl() bool {
	return sysctlPath() != ""
}

func NatEnabledIPv4() bool {
	return fileIs1("/proc/sys/net/ipv4/ip_forward")
}

func NatEnabledIPv6() bool {
	return fileIs1("/proc/sys/net/ipv6/conf/all/forwarding")
}

func EnableNat(ip4, ip6 bool) error {
	if ip4 {
		if err := writeLineToSysctlConf("net.ipv4.ip_forward=1"); err != nil {
			return err
		}
	}
	if ip6 {
		if err := writeLineToSysctlConf("net.ipv6.conf.all.forwarding=1"); err != nil {
			return err
		}
	}
	return reloadSysctl()
}

func writeLineToSysctlConf(line string) error {
	f, err := os.OpenFile("/etc/sysctl.conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open /etc/sysctl.conf: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(line + "\n")
	return err
}

func reloadSysctl() error {
	path := sysctlPath()
	if path == "" {
		return nil
	}
	cmd := exec.Command(path, "-p")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reload sysctl: %w", err)
	}
	return nil
}

func fileIs1(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	var buf [1]byte
	_, err = f.Read(buf[:])
	if err != nil {
		return false
	}
	return buf[0] == '1'
}

func sysctlPath() string {
	return binPath("sysctl")
}

func HasIPTables() bool {
	if os.Getenv("WGCMD_NO_DEPS") != "" {
		return true
	}
	return binPath("iptables") != ""
}

func pathNameForInterface(iface string) string {
	return "wgc-" + iface + ".path"
}

func serviceNameForInterface(iface string) string {
	return "wgc-" + iface + ".service"
}

func createSystemdPathForInterface(iface, wgdir string) error {
	wgConf := filepath.Join(wgdir, iface+".conf")
	pathContent := fmt.Sprintf("[Unit]\nDescription=Watch %s for changes", wgConf)
	pathContent += "\n\n[Path]\nPathModified=" + wgConf
	pathContent += "\n\n[Install]\nWantedBy=multi-user.target\n"
	pathName := pathNameForInterface(iface)
	pathLoc := filepath.Join("/etc/systemd/system", pathName)
	// save pathContent to pathLoc
	if err := os.WriteFile(pathLoc, []byte(pathContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", pathLoc, err)
	}

	// we don't need to reload systemd here, because we're going to do it later
	return nil
}

func createSystemdTargetService(iface string) error {
	pathName := pathNameForInterface(iface)
	targetServiceName := serviceNameForInterface(iface)

	targetService := fmt.Sprintf("[Unit]\nDescription=Restart WireGuard %s\nAfter=network.target", iface)
	targetService += fmt.Sprintf("\n\n[Service]\nType=oneshot\n"+
		"ExecStart=/usr/bin/systemctl restart wg-quick@%s.service", iface)
	targetService += "\n\n[Install]\nRequiredBy=" + pathName + "\n"

	targetServiceLoc := filepath.Join("/etc/systemd/system", targetServiceName)

	// save targetService to targetServiceLoc
	if err := os.WriteFile(targetServiceLoc, []byte(targetService), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetServiceLoc, err)
	}

	return nil
}

func reloadAndEnableSystemdStuff(iface string) error {
	pathName := pathNameForInterface(iface)
	targetServiceName := serviceNameForInterface(iface)

	// systemctl daemon-reload
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	// systemctl enable && start wg-quick@<iface>.service
	wgqSrv := "wg-quick@" + iface + ".service"
	cmd = exec.Command("systemctl", "enable", wgqSrv)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable %s: %w", wgqSrv, err)
	}

	cmd = exec.Command("systemctl", "start", wgqSrv)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start %s: %w", wgqSrv, err)
	}

	// systemctl enable wgc-<iface>-restart.service
	cmd = exec.Command("systemctl", "enable", pathName, targetServiceName)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		log.Println("cmd:", cmd.String())
		log.Println("err:", err)
		log.Println("stdout:", stdout.String())
		log.Println("stderr:", stderr.String())
		return fmt.Errorf("failed to enable %s,%s: %w", pathName, targetServiceName, err)
	}

	// systemctl start wgc-<iface>-restart.service
	cmd = exec.Command("systemctl", "start", pathName, targetServiceName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start %s,%s: %w", pathName, targetServiceName, err)
	}

	return nil
}

func CreateSystemdStuff(iface, wgdir string) error {
	if err := createSystemdPathForInterface(iface, wgdir); err != nil {
		return err
	}
	if err := createSystemdTargetService(iface); err != nil {
		return err
	}

	return reloadAndEnableSystemdStuff(iface)
}
