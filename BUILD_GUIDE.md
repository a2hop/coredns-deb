# CoreDNS with NAT64 and NAT664 Plugins - Build Guide

## Overview

This package builds CoreDNS from source with two custom plugins:

### NAT64 Plugin (Standard DNS64 - RFC 6147)
- Passes through IPv4 (A) queries normally
- Returns native AAAA records when they exist
- Only synthesizes IPv6 (AAAA) from IPv4 when no native AAAA exists
- Standard compliant DNS64 behavior

### NAT664 Plugin (Aggressive IPv6-only)
- Blocks IPv4 (A) record queries
- Always synthesizes IPv6 (AAAA) records from IPv4 addresses, even when native AAAA exists
- Forces IPv6-only operation

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
2. Integrate the nat64 and nat664 plugins
3. Build the CoreDNS binary with both plugins compiled in
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

Check that nat664 plugin is available:

```bash
coredns -plugins | grep nat664
```

Should output:
```
  dns.nat664
```

## Using the NAT664 Plugin

### Basic Configuration

Edit `/etc/coredns/Corefile`:

```
.:53 {
    errors
    log
    nat664 64:ff9b::
    forward . 8.8.8.8
    cache 30
}
```

The default prefix is `64:ff9b::` (RFC 6052 Well-Known Prefix). You can specify a custom prefix as shown above.

### Advanced Configuration

See `/etc/coredns/nat664.example` for a complete example.

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

The NAT664 plugin is located in `nat664/` directory:
- `nat664.go`: Main plugin logic
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
pkg/usr/local/bin/coredns -plugins | grep nat664
```

### Permission errors during build

The build script automatically fixes DEBIAN script permissions. If you still get errors:
```bash
chmod 755 pkg/DEBIAN/{postinst,prerm,postrm}
chmod 644 pkg/DEBIAN/control
```

## Development

To modify the NAT664 plugin:
1. Edit files in `nat664/` directory
2. Run `./build.sh` to rebuild
3. Test the new package

## References

- CoreDNS: https://coredns.io/
- NAT64 RFC: https://tools.ietf.org/html/rfc6146
- Go Proxy (China): https://goproxy.cn/
