#!/bin/bash
set -e

# Build script for CoreDNS Debian package with NAT64 plugin
# This script:
# 1. Clones CoreDNS source
# 2. Integrates the nat64 plugin
# 3. Builds CoreDNS binary
# 4. Creates the Debian package

VERSION=${1:-$(cat deb_version)}
COREDNS_VERSION=${2:-"v1.11.1"}
BUILD_DIR="build"
TEMP_DIR="${BUILD_DIR}/coredns-build"

echo "Building CoreDNS ${COREDNS_VERSION} with NAT64 plugin..."
echo "Package version: ${VERSION}"

# Set Go environment with China-friendly proxies
export GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct
export GOSUMDB=sum.golang.google.cn
export CGO_ENABLED=0

# Set git to use HTTPS instead of git protocol (better for proxies)
git config --global url."https://".insteadOf git://

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Clone CoreDNS
echo "Cloning CoreDNS..."
# Use GitHub mirror for China if available
GIT_URL="https://github.com/coredns/coredns.git"
# Uncomment the following line to use Gitee mirror (if you've mirrored it there)
# GIT_URL="https://gitee.com/mirrors/coredns.git"
git clone --depth 1 --branch ${COREDNS_VERSION} ${GIT_URL} ${TEMP_DIR}
cd ${TEMP_DIR}

# Copy NAT64 plugin
echo "Adding NAT64 plugin..."
mkdir -p plugin/nat64
cp ../../nat64/*.go plugin/nat64/

# Add nat64 to plugin.cfg (before 'forward' plugin for proper ordering)
echo "Configuring plugin..."
sed -i '/^forward:forward/i nat64:nat64' plugin.cfg

# Build CoreDNS
echo "Building CoreDNS..."
go generate
CGO_ENABLED=0 go build -ldflags="-s -w" -o coredns

# Verify the binary
echo "Verifying binary..."
./coredns -plugins | grep nat64 || (echo "ERROR: nat64 plugin not found in binary!" && exit 1)

# Copy binary to package structure
echo "Preparing package..."
cd ../..
mkdir -p pkg/usr/local/bin
cp ${TEMP_DIR}/coredns pkg/usr/local/bin/
chmod +x pkg/usr/local/bin/coredns

# Update package version
sed -i "s/^Version:.*/Version: ${VERSION}/" pkg/DEBIAN/control

# Build the deb package
echo "Building Debian package..."
PKG_NAME="coredns_${VERSION}_amd64"
mkdir -p ${BUILD_DIR}/${PKG_NAME}
cp -r pkg/* ${BUILD_DIR}/${PKG_NAME}/

# Fix permissions on DEBIAN scripts
chmod 755 ${BUILD_DIR}/${PKG_NAME}/DEBIAN/postinst
chmod 755 ${BUILD_DIR}/${PKG_NAME}/DEBIAN/prerm
chmod 755 ${BUILD_DIR}/${PKG_NAME}/DEBIAN/postrm
chmod 644 ${BUILD_DIR}/${PKG_NAME}/DEBIAN/control
chmod 644 ${BUILD_DIR}/${PKG_NAME}/DEBIAN/conffiles

dpkg-deb --build ${BUILD_DIR}/${PKG_NAME}

echo "Build complete!"
echo "Package: ${BUILD_DIR}/${PKG_NAME}.deb"
ls -lh ${BUILD_DIR}/${PKG_NAME}.deb

# Show plugin list
echo ""
echo "Installed plugins:"
pkg/usr/local/bin/coredns -plugins | grep -E "^(dns|nat64|forward|cache|errors|log)"
