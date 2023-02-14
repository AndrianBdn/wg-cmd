# WG Commander 
aka "wg-cmd" — TUI for managing WireGuard configuration files
![Screenshot](https://user-images.githubusercontent.com/994900/218503955-9f61eb8b-1eeb-4872-9f0e-979bed549d8d.png)

# Features
- no need for browser, forwarding HTTP ports - only SSH
- Setup Wizard
- TUI for managing peers
- view QR code in your Terminal
- automatically configures sysctl, systemd, NAT

## Current Limitations
- simple client-server Wireguard setup
- mostly for Linux 
- not a lot of managing besides adding / removing peers

# Installation 

Make sure you have Wireguard and iptables installed 
(`apt install wireguard-tools` in Ubuntu / Debian, `dnf install wireguard-tools iptables` in Rocky/Alma 9). 

Do download using curl execute (x86-64):
```
$ curl -SL https://github.com/andrianbdn/wg-cmd/releases/download/v0.1.0/wg-cmd-linux-amd64 -o /usr/local/bin/wg-cmd
```

Do download using curl execute (ARM 64):
```
$ curl -SL https://github.com/andrianbdn/wg-cmd/releases/download/v0.1.0/wg-cmd-linux-arm64 -o /usr/local/bin/wg-cmd
```

Set proper permissions and run the tool: 
```
$ chmod 755 /usr/local/bin/wg-cmd
$ wg-cmd
```

If you don't have /usr/local/bin in $PATH you will have to run `/usr/local/bin/wg-cmd` command. 

# Usage 

On first run WG Commander will show the setup wizard, that allows to configure new WireGuard interface interactively.

On subsequent runs (if wizard was successful) it will just display management TUI.

WG Commander requires root permissions to automatically tune sysctl, to create systemd units and to write to /etc/wireguard. 
You can avoid this if you know what you are doing. 

WG Commander keeps its own UI config in `~/.config/wg-cmd/wg-cmd.toml`

## Special options 

Run `wg-cmd new` to start the Wizard for new interface configuration

Run `wg-cmd <wg-interface>` to switch to specific interface (must be created before with wg-cmd)

Run `wg-cmd <wg-interface> make` to generate Wireguard configuration without showing UI.

# Tested 

WG Commander should work well on any systemd-based Linux distribution with wireguard, iptables, sysctl available.
It was tested on:

- Ubuntu 20.04
- Ubuntu 22.04
- Rocky Linux 9

# Notes 
There is no commercial purpose behind WG Commander. The project is licensed under the [MIT License](https://github.com/andrianbdn/wg-cmd/blob/master/LICENSE).

This project is NOT related to the creator of WireGuard®. WG Commander project is NOT approved, sponsored, or affiliated with WireGuard® or with the WireGuard® community.
