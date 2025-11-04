#!/bin/bash

# Verification script for CoreDNS NAT64 package
# This script checks if the package was built correctly

set -e

echo "=== CoreDNS NAT64 Package Verification ==="
echo ""

# Check if package exists
PKG_FILE="build/coredns_1.0.3_amd64.deb"
if [ ! -f "$PKG_FILE" ]; then
    echo "❌ Package not found: $PKG_FILE"
    echo "   Run ./build.sh first"
    exit 1
fi

echo "✓ Package found: $PKG_FILE"
ls -lh "$PKG_FILE"
echo ""

# Check binary exists in package
if dpkg -c "$PKG_FILE" | grep -q "./usr/local/bin/coredns$"; then
    echo "✓ CoreDNS binary found in package"
else
    echo "❌ CoreDNS binary not found in package"
    exit 1
fi

# Check NAT64 example config
if dpkg -c "$PKG_FILE" | grep -q "./etc/coredns/nat64.example"; then
    echo "✓ NAT64 example configuration found"
else
    echo "⚠  NAT64 example configuration not found"
fi

echo ""

# Check if binary has nat64 plugin
if [ -f "pkg/usr/local/bin/coredns" ]; then
    echo "Checking plugins in binary..."
    if pkg/usr/local/bin/coredns -plugins 2>/dev/null | grep -q "dns.nat64"; then
        echo "✓ NAT64 plugin is compiled into CoreDNS"
        echo ""
        echo "Available plugins:"
        pkg/usr/local/bin/coredns -plugins 2>/dev/null | grep -E "dns\.(nat64|forward|cache|errors|log|bind|hosts|file)"
    else
        echo "❌ NAT64 plugin NOT found in CoreDNS binary"
        exit 1
    fi
else
    echo "⚠  Binary not in pkg directory (normal after clean build)"
fi

echo ""
echo "=== Package Information ==="
dpkg-deb -I "$PKG_FILE" | grep -E "Package:|Version:|Architecture:|Description:"

echo ""
echo "=== Build Verification Complete ==="
echo ""
echo "To install:"
echo "  sudo dpkg -i $PKG_FILE"
echo ""
echo "After installation, verify with:"
echo "  coredns -plugins | grep nat64"
echo "  sudo systemctl status coredns"
