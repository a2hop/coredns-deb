#!/bin/bash

# Proxy configuration script for China
# Usage: source ./set-proxy.sh [proxy_url]
# Example: source ./set-proxy.sh http://127.0.0.1:7890

PROXY_URL="${1:-}"

if [ -z "$PROXY_URL" ]; then
    echo "Setting up China-friendly Go proxies (no HTTP proxy)..."
    export GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct
    export GOSUMDB=sum.golang.google.cn
    git config --global url."https://".insteadOf git://
    echo "✓ Go proxy set to: $GOPROXY"
else
    echo "Setting up HTTP/HTTPS proxy: $PROXY_URL"
    export HTTP_PROXY="$PROXY_URL"
    export HTTPS_PROXY="$PROXY_URL"
    export http_proxy="$PROXY_URL"
    export https_proxy="$PROXY_URL"
    export NO_PROXY="localhost,127.0.0.1,::1"
    export no_proxy="localhost,127.0.0.1,::1"
    
    # Also set Go proxies
    export GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct
    export GOSUMDB=sum.golang.google.cn
    
    # Git proxy
    git config --global http.proxy "$PROXY_URL"
    git config --global https.proxy "$PROXY_URL"
    git config --global url."https://".insteadOf git://
    
    echo "✓ HTTP/HTTPS proxy set to: $PROXY_URL"
    echo "✓ Go proxy set to: $GOPROXY"
fi

echo ""
echo "Current proxy settings:"
echo "  GOPROXY=$GOPROXY"
echo "  HTTP_PROXY=$HTTP_PROXY"
echo "  HTTPS_PROXY=$HTTPS_PROXY"
echo ""
echo "To unset proxy, run: source ./unset-proxy.sh"
