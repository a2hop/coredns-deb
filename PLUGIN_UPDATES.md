# Plugin Updates Summary

## Changes Made

### New NAT64 Plugin Added (Standard RFC 6147 DNS64)

A new `nat64` plugin has been added alongside the existing `nat664` plugin. Both plugins are now included in the CoreDNS build.

### Key Differences Between the Plugins

#### NAT64 (Standard DNS64 - RFC 6147)
- **Standards compliant** DNS64 implementation
- Passes through A queries unchanged
- Returns native AAAA records when they exist
- Only synthesizes AAAA from A records when no native AAAA exists
- Ideal for dual-stack networks where native IPv6 should be preferred

**Key Code Feature:**
```go
// Check if native AAAA records exist
if len(aaaaResp.Answer) > 0 {
    for _, rr := range aaaaResp.Answer {
        if rr.Header().Rrtype == dns.TypeAAAA {
            // Native AAAA exists, don't synthesize - return as-is
            w.WriteMsg(aaaaResp)
            return dns.RcodeSuccess, nil
        }
    }
}
```

#### NAT664 (Aggressive IPv6-only)
- **Custom implementation** for forcing IPv6-only operation
- Blocks all A queries with NXDOMAIN
- Always synthesizes AAAA records, even when native AAAA exists
- Forces all traffic through NAT64 gateway
- Ideal for pure IPv6-only networks

### Files Created/Modified

#### New Files Created:
1. `/nat64/nat64.go` - Main plugin implementation (standard DNS64)
2. `/nat64/setup.go` - Plugin setup and configuration
3. `/nat64/go.mod` - Go module definition
4. `/nat64/README.md` - Plugin documentation
5. `/pkg/etc/coredns/nat64.example` - Example configuration

#### Modified Files:
1. `build.sh` - Updated to build both nat64 and nat664 plugins
2. `BUILD_GUIDE.md` - Updated documentation for both plugins
3. `README.md` - Updated feature list and examples for both plugins
4. `pkg/DEBIAN/control` - Updated package description
5. `.github/workflows/build_and_release.yaml` - Updated CI/CD to include both plugins

### Configuration Examples

#### For Standard DNS64 (nat64):
```corefile
.:53 {
    errors
    log
    nat64 64:ff9b::
    forward . 8.8.8.8
    cache 30
}
```

#### For Aggressive IPv6-only (nat664):
```corefile
.:53 {
    errors
    log
    nat664 64:ff9b::
    forward . 8.8.8.8
    cache 30
}
```

### Use Cases

**Use NAT64 when:**
- You have a dual-stack network
- You want RFC-compliant DNS64 behavior
- You want to prefer native IPv6 when available
- You need A queries to work normally

**Use NAT664 when:**
- You have an IPv6-only network
- You want to force all traffic through NAT64
- You want to block all IPv4 queries
- You need guaranteed IPv6 operation

### Build Process

The build script now:
1. Creates `plugin/nat64/` directory
2. Copies nat64 plugin files
3. Creates `plugin/nat664/` directory
4. Copies nat664 plugin files
5. Adds both plugins to CoreDNS `plugin.cfg`
6. Builds CoreDNS with both plugins
7. Verifies both plugins are present

### Testing

Both plugins can be tested in the same installation:

```bash
# Build the package
./build.sh

# Install
sudo dpkg -i build/coredns_*.deb

# Verify both plugins are available
coredns -plugins | grep nat64
coredns -plugins | grep nat664

# Test with nat64 configuration
sudo cp /etc/coredns/nat64.example /etc/coredns/Corefile
sudo systemctl restart coredns
dig @localhost AAAA google.com  # Should return native AAAA if exists
dig @localhost A google.com     # Should work normally

# Test with nat664 configuration
sudo cp /etc/coredns/nat664.example /etc/coredns/Corefile
sudo systemctl restart coredns
dig @localhost AAAA google.com  # Should return synthesized AAAA
dig @localhost A google.com     # Should return NXDOMAIN
```

### Documentation

- `/nat64/README.md` - Detailed NAT64 plugin documentation
- `/nat664/README.md` - Detailed NAT664 plugin documentation
- `/pkg/etc/coredns/nat64.example` - NAT64 configuration example
- `/pkg/etc/coredns/nat664.example` - NAT664 configuration example
- `README.md` - Updated main documentation
- `BUILD_GUIDE.md` - Updated build instructions

### Next Steps

1. Test the build process: `./build.sh`
2. Verify both plugins are compiled: `pkg/usr/local/bin/coredns -plugins`
3. Test both configurations with real DNS queries
4. Update version in `deb_version` if needed
5. Commit and push changes to trigger CI/CD build
