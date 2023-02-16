# WG Commander 

aka "wg-cmd" — TUI for managing WireGuard configuration files

<a href="https://github.com/andrianbdn/wg-cmd/releases/latest"><img src="https://img.shields.io/github/v/release/andrianbdn/wg-cmd" /></a>
<a href="./LICENSE"><img src="https://img.shields.io/github/license/andrianbdn/wg-cmd" /></a>
<a href="./go.mod"><img src="https://img.shields.io/github/go-mod/go-version/andrianbdn/wg-cmd" /></a>
[![Go Report Card](https://goreportcard.com/badge/github.com/andrianbdn/wg-cmd)](https://goreportcard.com/report/github.com/andrianbdn/wg-cmd)

![screenshot](https://user-images.githubusercontent.com/994900/218720566-e5b3ab22-d7fc-4df7-a777-ad9b6280ada8.png)

# Features
- no need for browser, forwarding HTTP ports - only SSH
- has a nice Setup Wizard
- Text-based user interface for managing peers
- view QR code in your Terminal
- automatically configures sysctl, systemd, NAT

## Current Limitations
- only simple client-server WireGuard setup
- mostly for Linux (assumes iptables, systemd, sysctl are available) — see [Other OS](#other-os-besides-linux) section
- can't manage existing Wireguard interfaces (but you can create new WireGuard interfaces on the same host)

# Installation 

Make sure you have Wireguard and iptables installed 
(`apt install wireguard-tools` in Ubuntu / Debian, `dnf install wireguard-tools iptables` in Rocky/Alma 9). 

To download using curl run:
```shell
# for x86_64 this command 
curl -SL https://github.com/andrianbdn/wg-cmd/releases/download/v0.1.1/wg-cmd-0.1.1-linux-amd64 -o /usr/local/bin/wg-cmd

# for arm64 this command
curl -SL https://github.com/andrianbdn/wg-cmd/releases/download/v0.1.1/wg-cmd-0.1.1-linux-arm64 -o /usr/local/bin/wg-cmd
```

Set proper permissions and run the tool: 
```
$ chmod 755 /usr/local/bin/wg-cmd
$ wg-cmd
```

If you don't have /usr/local/bin in $PATH you will have to
run `/usr/local/bin/wg-cmd` command using full path.

# Usage 

On first run WG Commander will show the setup wizard, that allows to configure new WireGuard interface interactively.

On subsequent runs (if wizard was successful) it will just display management TUI.

## Advanced usage

WG Commander requires root permissions to automatically tune sysctl, to create systemd units and to write to /etc/wireguard. 
You can avoid this if you know what you are doing. 

WG Commander keeps its own UI config in `~/.config/wg-cmd/wg-cmd.toml`

### Special options 

Run `wg-cmd new` to start the Wizard for new interface configuration

Run `wg-cmd <wg-interface>` to switch to specific interface (must be created before with wg-cmd)

Run `wg-cmd <wg-interface> make` to generate Wireguard configuration without showing UI.

### Configuration 

WG Commander uses directories as a "database". 
It stores the interface configuration in /etc/wireguard/wgc-<interface-name> directory. 

The configuration is stored in [TOML](https://toml.io) files.

Most configuration keys are similar to WireGuard ones. 

#### server configuration (0001-server.toml)
Some keys in this configuration file will actually be used for generating 
client configuration files. 

`ClientRoute` - AllowedIPs for client config

`ClientDNS` - DNS for client config

`ClientServerEndpoint` - Endpoint for client config

`ClientPersistentKeepalive` - PersistentKeepalive for client config 

#### client configuration (nnnn-%client%.toml)

Contains `ClientRoute` that overrides the one from server config.

Client files contains `PrivateKey` field. 
If you find it unacceptable, you can remove it from the file after you exported 
configuration (or QR code) to the client.

### Other OS besides Linux

WG Commander is designed to work on Linux, because it uses systemd, iptables, sysctl. 
However, it is written in plain Go, so it should work on any OS that Go supports.

- You will need to compile binary yourself.
- Set the environment variable `WG_CMD_NO_DEPS` to 1 to disable any Linux-specific checking on start. 
- Edit 0001-server.toml and set your OS commands in PostUp4/PostUp6/PostDown4/PostDown6 fields.
- You will need to arrange WireGuard configuration reload yourself: wither restart WireGuard manually
when config changes, or monitor the /etc/wireguard/wg*.conf files for changes and reload automatically.

PRs are welcome to add support for other OSes.

### Uninstall 

To uninstall WG Commander, just remove the binary from /usr/local/bin/wg-cmd. 
You can also remove directories /etc/wireguard/wgc-* and ~/.config/wg-cmd

If you have created systemd units, you will need to remove them manually.

Below is an example of how to remove WG Commander managed interface wg7 
(change it to whatever interface you need to delete):

```sh
systemctl stop wgc-wg7.{path,service}
systemctl disable wgc-wg7.{path,service}
rm /etc/systemd/system/wgc-wg7.{path,service}
systemctl stop wg-quick@wg7.service
systemctl disable wg-quick@wg7.service
rm /etc/wireguard/wg7.conf
rm -Rf /etc/wireguard/wgc-wg7
```


# Tested
WG Commander should work well on any systemd-based Linux
distribution with WireGuard, iptables, sysctl available.
It was tested on:
- Ubuntu 20.04
- Ubuntu 22.04
- Rocky Linux 9
- Debian 11

# Notes 
There is no commercial purpose behind WG Commander. 
The project is licensed under 
the [MIT License](https://github.com/andrianbdn/wg-cmd/blob/master/LICENSE).

This project is NOT related to the creator of WireGuard®.
WG Commander project is NOT approved, sponsored, or affiliated 
with WireGuard® or with the WireGuard® community.
