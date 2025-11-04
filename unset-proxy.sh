#!/bin/bash

# Unset all proxy settings

echo "Unsetting proxy configuration..."

unset HTTP_PROXY
unset HTTPS_PROXY
unset http_proxy
unset https_proxy
unset NO_PROXY
unset no_proxy

git config --global --unset http.proxy
git config --global --unset https.proxy

# Keep Go proxy for China
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

echo "✓ HTTP/HTTPS proxy unset"
echo "✓ Go proxy still set to: $GOPROXY (recommended for China)"
