# NAT64 Plugin for CoreDNS

This is a custom CoreDNS plugin that provides NAT64 DNS64 functionality.

## What it does

The NAT64 plugin:
- **Blocks IPv4 (A) queries**: Returns NXDOMAIN for A record requests
- **Synthesizes IPv6 (AAAA) records**: Converts IPv4 addresses to IPv6 using a configured prefix
- **Passes through other queries**: All other DNS record types are forwarded normally

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

## How it works

1. When a client queries for an AAAA record (IPv6):
   - The plugin queries upstream for the corresponding A record (IPv4)
   - Takes the IPv4 address and embeds it in the configured IPv6 prefix
   - Returns the synthesized AAAA record to the client

2. When a client queries for an A record (IPv4):
   - Returns NXDOMAIN (name error)

3. All other query types (MX, TXT, etc.) are passed through unchanged

## Example

Query: `AAAA google.com`
- Upstream returns: `A 142.250.185.46`
- Plugin synthesizes: `AAAA 64:ff9b::8efa:b92e`
  - IPv4 `142.250.185.46` = hex `8e.fa.b9.2e`
  - Embedded in RFC 6052 prefix: `64:ff9b::8efa:b92e`

With custom prefix `2001:db8:64::`:
- Same IPv4 `142.250.185.46`
- Result: `AAAA 2001:db8:64::8efa:b92e`

## Testing

```bash
# Test AAAA query (should return synthesized IPv6)
dig @localhost AAAA google.com

# Test A query (should return NXDOMAIN)
dig @localhost A google.com

# Test other types (should work normally)
dig @localhost MX google.com
```

## Installation

This plugin is integrated into the CoreDNS binary during the package build process.

See the main README for build instructions.
