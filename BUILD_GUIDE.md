# CoreDNS with NAT64 Plugin - Build Guide

## Overview

This package builds CoreDNS from source with a custom NAT64 plugin that:
- Blocks IPv4 (A) record queries
- Synthesizes IPv6 (AAAA) records from IPv4 addresses using a configurable NAT64 prefix
- Passes through all other query types

## Building the Package

### Prerequisites

```bash
sudo apt install -y golang-go git dpkg-dev
```

### Quick Build

```bash
./build.sh
```

This will:
1. Clone CoreDNS source (v1.11.1)
2. Integrate the nat64 plugin
3. Build the CoreDNS binary with the plugin compiled in
4. Create a Debian package: `build/coredns_1.0.3_amd64.deb`

### Build with Custom Version

```bash
./build.sh 1.0.4 v1.11.3
```

Arguments:
- First arg: Package version (default: from `deb_version` file)
- Second arg: CoreDNS git tag (default: v1.11.1)

### For China Users

The build script automatically uses China-friendly proxies:
- GOPROXY: goproxy.cn, mirrors.aliyun.com/goproxy/, goproxy.io
- GOSUMDB: sum.golang.google.cn

If you need additional HTTP/HTTPS proxy:

```bash
source ./set-proxy.sh http://127.0.0.1:7890
./build.sh
source ./unset-proxy.sh
```

## Installing the Package

```bash
sudo dpkg -i build/coredns_1.0.3_amd64.deb
```

## Verifying the Installation

Check that nat64 plugin is available:

```bash
coredns -plugins | grep nat64
```

Should output:
```
  dns.nat64
```

## Using the NAT64 Plugin

### Basic Configuration

Edit `/etc/coredns/Corefile`:

```
.:53 {
    errors
    log
    nat64 64:ff9b::
    forward . 8.8.8.8
    cache 30
}
```

The default prefix is `64:ff9b::` (RFC 6052 Well-Known Prefix). You can specify a custom prefix as shown above.

### Advanced Configuration

See `/etc/coredns/nat64.example` for a complete example.

### Starting CoreDNS

```bash
sudo systemctl start coredns
sudo systemctl enable coredns
sudo systemctl status coredns
```

## Testing

Test IPv4 blocking:
```bash
dig @127.0.0.1 google.com A
# Should return NXDOMAIN
```

Test IPv6 synthesis:
```bash
dig @127.0.0.1 google.com AAAA
# Should return synthesized IPv6 addresses
```

## Plugin Details

The NAT64 plugin is located in `nat64/` directory:
- `nat64.go`: Main plugin logic
- `setup.go`: Plugin registration and configuration
- `README.md`: Plugin documentation

## Package Contents

- **Binary**: `/usr/local/bin/coredns` (54MB+)
- **Config**: `/etc/coredns/Corefile`
- **Examples**: `/etc/coredns/*.example`
- **Service**: `/etc/systemd/system/coredns.service`
- **Scripts**: `/usr/local/bin/coredns-*` (helper utilities)

## Troubleshooting

### Build fails to clone CoreDNS

Check your internet connection or use a proxy:
```bash
source ./set-proxy.sh http://your-proxy:port
```

### Plugin not showing up

Verify the binary was built correctly:
```bash
pkg/usr/local/bin/coredns -plugins | grep nat64
```

### Permission errors during build

The build script automatically fixes DEBIAN script permissions. If you still get errors:
```bash
chmod 755 pkg/DEBIAN/{postinst,prerm,postrm}
chmod 644 pkg/DEBIAN/control
```

## Development

To modify the NAT64 plugin:
1. Edit files in `nat64/` directory
2. Run `./build.sh` to rebuild
3. Test the new package

## References

- CoreDNS: https://coredns.io/
- NAT64 RFC: https://tools.ietf.org/html/rfc6146
- Go Proxy (China): https://goproxy.cn/
