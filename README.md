# WG Commander 

aka "wg-cmd" — TUI for managing WireGuard configuration files

<a href="https://github.com/andrianbdn/wg-cmd/releases/latest"><img src="https://img.shields.io/github/v/release/andrianbdn/wg-cmd" /></a>
<a href="./LICENSE"><img src="https://img.shields.io/github/license/andrianbdn/wg-cmd" /></a>
<a href="./go.mod"><img src="https://img.shields.io/github/go-mod/go-version/andrianbdn/wg-cmd" /></a>
[![Go Report Card](https://goreportcard.com/badge/github.com/andrianbdn/wg-cmd)](https://goreportcard.com/report/github.com/andrianbdn/wg-cmd)

![screenshot](https://user-images.githubusercontent.com/994900/218720566-e5b3ab22-d7fc-4df7-a777-ad9b6280ada8.png)

# Features
- no need for a browser or HTTP port - works in the terminal, over SSH too
- has a nice Setup Wizard
- text-based user interface for managing peers
- view QR code in the terminal
- automatically configures sysctl, systemd, NAT

## Current Limitations
- supports simple client-server WireGuard setup
- mostly for Linux (assumes iptables, systemd, sysctl are available) — see [Other OS](#other-os-besides-linux) section
- can't manage existing WireGuard configuration (but you can create new WireGuard interfaces on the same host)

# Installation 

Make sure you have WireGuard and iptables installed 
(`apt install wireguard-tools` in Ubuntu / Debian, `dnf install wireguard-tools iptables` in Rocky/Alma 9). 

To download using curl run:
```shell
# for x86_64 
curl -SL https://github.com/andrianbdn/wg-cmd/releases/download/v0.1.6/wg-cmd-0.1.6-linux-amd64 -o /usr/local/bin/wg-cmd

# for arm64 
curl -SL https://github.com/andrianbdn/wg-cmd/releases/download/v0.1.6/wg-cmd-0.1.6-linux-arm64 -o /usr/local/bin/wg-cmd
```

Set proper permissions and run the tool: 
```
chmod 755 /usr/local/bin/wg-cmd
wg-cmd
```

If you don't have `/usr/local/bin` in $PATH you will have to
run `/usr/local/bin/wg-cmd` command using the full path.

WG Commander requires root permissions to automatically tune sysctl, to create systemd units and to write to /etc/wireguard.

# Usage 

On first run WG Commander will show the setup wizard, that allows to configure new WireGuard interface interactively.

On subsequent runs (if wizard was successful) it will just display management TUI.

Note regarding the QR code: some devices (Android?) may require a higher quality QR code. WG Commander will automatically increase quality when you make Terminal window resolution larger (smaller font, larger window).

## Advanced usage

You can run WG Commander as a non-root user if you change permissions on 
/etc/wireguard and configure sysctl/systemd manually.

WG Commander keeps its own UI config in `~/.config/wg-cmd/wg-cmd.toml`

The most important options are:
```toml
WireguardDir = "/etc/wireguard"
# directory for WireGuard configuration files 

DatabaseDir = "/etc/wireguard"
# directory for WG Commander database files (wgc-<interface-name>)
```

You can change these options to point to directories that you have write access to.

### Special options 

Run `wg-cmd new` to start the wizard for new interface configuration.

Run `wg-cmd <wg-interface>` to switch to specific interface (must be created before with wg-cmd).

Run `wg-cmd <wg-interface> make` to generate Wireguard configuration without showing UI.

### Configuration 

WG Commander uses directories as its "database". 
It stores the interface configuration in /etc/wireguard/wgc-<interface-name> directory. 

The configuration is stored using [TOML](https://toml.io) file format.

Most configuration keys are similar to WireGuard ones. 

#### server configuration (0001-server.toml)
Some keys in this configuration file will actually be used for generating 
client configuration files. 

`ClientRoute` - AllowedIPs for client config

`ClientDNS` - DNS configuration value for all clients

`ClientServerEndpoint` - Endpoint for client config

`ClientPersistentKeepalive` - PersistentKeepalive for client config 

`MTU` - MTU for the server and client (0 — make WireGuard choose)

#### client configuration (nnnn-%client%.toml)

`ClientRoute` - completely overrides the `ClientRoute` from the server config

`AddServerRoute` - adds additional network to AllowedIPs for the client on the server side (useful 
when you want to route traffic to one client to another client's network through the server)

`MTU` - Override server MTU with a different value for this client. Set to -1 to omit MTU from this WireGuard client config.

`DNS` - Override server `ClientDNS` setting for all clients. Specify a comma separated IP list. 
Set to `no` or `none` to omit DNS from this WireGuard client config.

Client configuration files contain `PrivateKey` field. 
If you find it unacceptable, you can remove it from the file after you exported 
configuration (or QR code) to the client.

### Other OS besides Linux

WG Commander is designed to work on Linux, because it uses procfs, systemd, iptables, sysctl. 
However, it is written in plain Go, so it should work on any OS that Go supports.

- You will need to compile binary yourself.
- Set the environment variable `WG_CMD_NO_DEPS` to 1 to disable any Linux-specific checks on start. 
- Edit 0001-server.toml and set your OS commands in PostUp4/PostUp6/PostDown4/PostDown6 fields.
- You will need to reload WireGuard configuration: manually when you change something
or monitor /etc/wireguard/wg*.conf files for changes and reload WireGuard automatically.

PRs are welcome to add support for other OSes.

### Running in Docker 

Although it is possible, it is not recommended to run WG Commander in Docker. 

The Setup Wizard will not work properly, because it needs to create systemd units and modify sysctl.


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
distribution with WireGuard, iptables, sysctl, procfs available.

It was tested on:
- Ubuntu 24.04 (v0.1.6 tested on Aug 17 2024)
- Ubuntu 20.04
- Ubuntu 22.04
- Rocky Linux 9
- Debian 11
- Debian 12

# Notes 
There is no commercial purpose behind WG Commander. 
The project is licensed under 
the [MIT License](https://github.com/andrianbdn/wg-cmd/blob/master/LICENSE).

This project is NOT related to the creator of WireGuard®.
WG Commander project is NOT approved, sponsored, or affiliated 
with WireGuard® or with the WireGuard® community.
