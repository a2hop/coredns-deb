# CoreDNS Debian Package

Complete Debian package for CoreDNS with systemd integration, intelligent upgrade handling, and port 53 conflict management.

## Quick Start

### Installation
```bash
sudo apt install ./coredns-1.0.0-amd64.deb
```

### Basic Setup
```bash
# 1. Resolve port 53 conflict with systemd-resolved
sudo coredns-resolve disable

# 2. Configure CoreDNS
sudo nano /etc/coredns/Corefile

# 3. Start service
sudo systemctl start coredns
sudo systemctl enable coredns

# 4. Test DNS
dig @127.0.0.1
```

## Key Features

- ✅ **Intelligent Upgrades**: Only restarts if service was running before upgrade
- ✅ **systemd-resolved Integration**: One-command to resolve port 53 conflicts
- ✅ **CLI Tools**: `coredns-ctl`, `coredns-zone`, `coredns-validate`, `coredns-resolve`
- ✅ **Shell Completion**: Tab completion for all utilities
- ✅ **Security**: AppArmor profile, unprivileged user execution
- ✅ **Documentation**: Man pages, examples, comprehensive guides

## Managing systemd-resolved Conflicts

### Problem
Both CoreDNS and systemd-resolved try to use port 53.

### Solution
Use `coredns-resolve`:

```bash
# Check current status
coredns-resolve status

# Use CoreDNS (disable systemd-resolved stub listener)
sudo coredns-resolve disable

# Revert to systemd-resolved
sudo coredns-resolve enable

# Interactive configuration
sudo coredns-resolve configure

# Check what's on port 53
coredns-resolve check
```

Via `coredns-ctl`:
```bash
sudo coredns-ctl resolve status
sudo coredns-ctl resolve disable
```

## CLI Tools

### coredns-ctl - Service Control
```bash
coredns-ctl start              # Start service
coredns-ctl stop               # Stop service
coredns-ctl restart            # Restart service
coredns-ctl status             # Service status
coredns-ctl logs               # Show logs (last 50 lines)
coredns-ctl follow             # Follow logs in real-time
coredns-ctl config             # Show configuration
coredns-ctl validate           # Validate configuration
coredns-ctl resolve status     # Resolve subcommand
```

### coredns-zone - Zone Management
```bash
coredns-zone init example.com          # Create new zone
coredns-zone bump example.com          # Increment serial number
coredns-zone list                      # List zones
coredns-zone show example.com          # Show zone file
```

### coredns-validate - Configuration Validator
```bash
coredns-validate /etc/coredns/Corefile
```

### coredns-resolve - systemd-resolved Management
```bash
coredns-resolve status      # Show current configuration
coredns-resolve disable     # Free port 53 for CoreDNS
coredns-resolve enable      # Restore systemd-resolved
coredns-resolve check       # Check port 53 usage
coredns-resolve configure   # Interactive wizard
```

## Upgrade Behavior

The package intelligently handles upgrades:

- **Service Running**: Automatically restarts with new version
- **Service Stopped**: Remains stopped, respects your intent

Example:
```bash
# Service is running
sudo systemctl start coredns

# Upgrade CoreDNS
sudo apt upgrade coredns
# Output: "CoreDNS service restarted successfully"

# Service was stopped
sudo systemctl stop coredns

# Upgrade CoreDNS
sudo apt upgrade coredns
# Output: "Service remains stopped, start when ready: sudo systemctl start coredns"
```

## Configuration

### Default Configuration
`/etc/coredns/Corefile` - Created during installation

### Example Configurations
```bash
# View examples
ls /etc/coredns/*.example

# Copy an example to use
sudo cp /etc/coredns/proxy_zone.example /etc/coredns/Corefile

# Validate before restarting
coredns-validate /etc/coredns/Corefile

# Reload without downtime
sudo coredns-ctl reload
```

### Zone Files
```bash
# Create example zone
sudo coredns-zone init example.com

# Zone file created at: /etc/coredns/example.com.zone
# Reference in Corefile:
# example.com {
#     file /etc/coredns/example.com.zone
# }
```

## Example Configurations Included

| File | Purpose |
|------|---------|
| `proxy_only.example` | Simple DNS proxy |
| `proxy_zone.example` | Local zone + proxy |
| `proxy_zone_conditional.example` | Conditional forwarding with health checks |
| `test.example` | Test configuration |

## Systemd Service Management

```bash
# View service status
systemctl status coredns

# Enable auto-start
sudo systemctl enable coredns

# Disable auto-start
sudo systemctl disable coredns

# View logs
sudo journalctl -u coredns

# Follow logs real-time
sudo journalctl -u coredns -f

# View recent errors
sudo journalctl -u coredns -n 50
```

## Logging

CoreDNS logs to systemd journal:

```bash
# View all CoreDNS logs
sudo journalctl -u coredns

# Follow in real-time
sudo journalctl -u coredns -f

# Last 50 entries
sudo journalctl -u coredns -n 50

# Clean logs older than 7 days
sudo journalctl --vacuum=time=7d

# Or use coredns-ctl
coredns-ctl logs       # Last 50 lines
coredns-ctl follow     # Real-time follow
```

## File Locations

```
/usr/local/bin/
  coredns              # Main binary
  coredns-ctl          # Service control utility
  coredns-zone         # Zone management
  coredns-validate     # Configuration validator
  coredns-resolve      # systemd-resolved management

/etc/coredns/
  Corefile             # Main configuration
  *.example            # Example configurations
  *.zone               # Zone files

/etc/systemd/system/
  coredns.service      # Systemd service file

/usr/share/man/man1/
  coredns-ctl.1        # Man page for coredns-ctl
  coredns-zone.1       # Man page for coredns-zone
  coredns-resolve.1    # Man page for coredns-resolve
```

## Troubleshooting

### Port 53 Already in Use
```bash
# Check what's using it
coredns-resolve check

# Disable systemd-resolved
sudo coredns-resolve disable

# Or use interactive wizard
sudo coredns-resolve configure
```

### Service Won't Start
```bash
# Check status and errors
sudo systemctl status coredns

# View detailed logs
sudo journalctl -u coredns -n 100

# Validate configuration
coredns-validate /etc/coredns/Corefile

# Test manually
sudo /usr/local/bin/coredns -conf=/etc/coredns/Corefile
```

### DNS Resolution Not Working
```bash
# Check service is running
sudo systemctl status coredns

# Test local DNS
dig @127.0.0.1

# Check for errors
coredns-ctl logs

# Verify configuration
coredns-ctl config
```

## Security Features

- **Unprivileged User**: Runs as unprivileged `coredns` user
- **Capability Limiting**: Only `CAP_NET_BIND_SERVICE` capability
- **AppArmor Profile**: Optional security profile available
- **Secure Permissions**: Sensible defaults for all files
- **Configuration Protection**: Debian conffiles protect user configs

## Installation & Removal

### Install
```bash
sudo apt install ./coredns-1.0.0-amd64.deb
```

### Remove (keep config)
```bash
sudo apt remove coredns
```

### Remove (complete cleanup)
```bash
sudo apt purge coredns
```

## Package Contents

- **5 CLI Utilities**: Control, zone management, validation, resolve
- **4 System Scripts**: Installation, pre-removal, post-removal
- **8 Configuration Files**: Default config, examples, systemd files
- **Documentation**: README, man pages, copyright
- **Security**: AppArmor profile, bash completion, logrotate

## Building the Package

### Prerequisites
```bash
# On Debian/Ubuntu
sudo apt install build-essential devscripts dpkg-dev

# Download CoreDNS binary
wget https://github.com/coredns/coredns/releases/download/v1.10.1/coredns_linux_amd64.tgz
tar xzf coredns_linux_amd64.tgz
```

### Build
```bash
# Place binary in package
cp coredns pkg/usr/local/bin/

# Build package
dpkg-deb --build pkg coredns-1.0.0-amd64.deb

# Verify
dpkg -c coredns-1.0.0-amd64.deb
```

## Help & Documentation

```bash
# View help for any utility
coredns-ctl --help
coredns-zone --help
coredns-validate --help
coredns-resolve --help

# Read man pages
man coredns-ctl
man coredns-zone
man coredns-resolve

# View installed configuration
cat /etc/coredns/Corefile

# List available examples
ls -la /etc/coredns/*.example
```

## Common Use Cases

### Use Case 1: DNS Proxy
```bash
# Copy proxy example
sudo cp /etc/coredns/proxy_only.example /etc/coredns/Corefile

# Edit to set upstream DNS
sudo nano /etc/coredns/Corefile

# Restart service
sudo coredns-ctl restart
```

### Use Case 2: Local Domain + Proxy
```bash
# Copy example
sudo cp /etc/coredns/proxy_zone.example /etc/coredns/Corefile

# Create zone
sudo coredns-zone init internal.company.com

# Edit zone file
sudo nano /etc/coredns/internal.company.com.zone

# Restart service
sudo coredns-ctl restart
```

### Use Case 3: Split Horizon DNS
```bash
# Use conditional forwarding example
sudo cp /etc/coredns/proxy_zone_conditional.example /etc/coredns/Corefile

# Customize for your networks
sudo nano /etc/coredns/Corefile

# Validate and restart
coredns-validate /etc/coredns/Corefile
sudo coredns-ctl restart
```

## Support

For issues or questions:
- Check logs: `coredns-ctl logs`
- Validate config: `coredns-validate /etc/coredns/Corefile`
- View man pages: `man coredns-ctl`
- Check CoreDNS docs: https://coredns.io/

## License

See `/usr/share/doc/coredns/copyright` for license information.
