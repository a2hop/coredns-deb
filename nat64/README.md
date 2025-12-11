# NAT64 Plugin for CoreDNS (RFC 6147 DNS64)

This is a standard DNS64 implementation following RFC 6147. It provides NAT64/DNS64 functionality for IPv6-only networks.

## What it does

The NAT64 plugin follows the standard DNS64 behavior:
- **Passes through A queries**: IPv4 (A) queries are forwarded normally
- **Returns native AAAA records**: If real IPv6 addresses exist, they are returned unchanged
- **Synthesizes AAAA only when needed**: Converts IPv4 addresses to IPv6 only when no native AAAA records exist
- **Passes through other queries**: All other DNS record types are forwarded normally

## Difference from NAT664

This package includes two plugins:
- **nat64**: Standard RFC 6147 DNS64 (this plugin) - synthesizes only when no native AAAA exists
- **nat664**: Aggressive NAT64 - always blocks A queries and always synthesizes, even when native AAAA exists

Use **nat64** for standard DNS64 behavior compatible with dual-stack networks.
Use **nat664** to force IPv6-only operation.

## Configuration

### Syntax

```
nat64 [PREFIX]
```

- **PREFIX**: IPv6 prefix to use for synthesis (default: `64:ff9b::` - RFC 6052 Well-Known Prefix)

### Examples

#### Basic usage with default prefix

```corefile
.:53 {
    errors
    log
    nat64
    forward . 8.8.8.8
    cache 30
}
```

#### Custom NAT64 prefix

```corefile
.:53 {
    errors
    log
    nat64 2001:db8:64::
    forward . 8.8.8.8
    cache 30
}
```

#### Zone-specific NAT64

```corefile
example.com {
    nat64 2001:db8:64::
    forward . 192.168.1.1
}

. {
    forward . 8.8.8.8
}
```

## How it works (Standard DNS64 - RFC 6147)

1. When a client queries for an AAAA record (IPv6):
   - The plugin first checks if native AAAA records exist
   - **If native AAAA exists**: Returns them unchanged (no synthesis)
   - **If no AAAA exists**: Queries for A records and synthesizes AAAA from them
   
2. When a client queries for an A record (IPv4):
   - Passes through normally (standard behavior)

3. All other query types (MX, TXT, etc.) are passed through unchanged

## Example

### Case 1: Domain has native IPv6
Query: `AAAA google.com`
- Upstream returns: `AAAA 2607:f8b0:4004:c07::64`
- Plugin returns: `AAAA 2607:f8b0:4004:c07::64` (no synthesis)

### Case 2: Domain has only IPv4
Query: `AAAA example-ipv4-only.com`
- Upstream has no AAAA, but has: `A 192.0.2.1`
- Plugin synthesizes: `AAAA 64:ff9b::c000:201`
  - IPv4 `192.0.2.1` = hex `c0.00.02.01`
  - Embedded in RFC 6052 prefix: `64:ff9b::c000:201`

With custom prefix `2001:db8:64::`:
- Same IPv4 `192.0.2.1`
- Result: `AAAA 2001:db8:64::c000:201`

## Testing

```bash
# Test AAAA query for domain with native IPv6 (should return native AAAA)
dig @localhost AAAA google.com

# Test AAAA query for IPv4-only domain (should return synthesized AAAA)
dig @localhost AAAA ipv4-only-example.com

# Test A query (should work normally)
dig @localhost A google.com

# Test other types (should work normally)
dig @localhost MX google.com
```

## Standards Compliance

This plugin implements RFC 6147 (DNS64) with the following features:
- ✅ Returns native AAAA records when they exist
- ✅ Synthesizes AAAA from A only when no AAAA exists
- ✅ Passes through A queries unchanged
- ✅ Uses RFC 6052 Well-Known Prefix (64:ff9b::/96) by default
- ✅ Correctly embeds IPv4 in last 32 bits of IPv6 address
- ✅ Preserves CNAME chains
- ✅ Filters DNSSEC records to prevent validation failures

## Installation

This plugin is integrated into the CoreDNS binary during the package build process.
