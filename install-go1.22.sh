#!/bin/bash
# 方案二：在用户目录安装 Go 1.22，供 Cursor 安装 gopls 等工具使用
# 运行: bash install-go1.22.sh

set -e
GO_VERSION=1.22.13
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/${GO_TAR}"
INSTALL_DIR="$HOME/go1.22"

echo "下载 Go ${GO_VERSION}..."
cd /tmp
rm -f "$GO_TAR"
wget "$GO_URL" || curl -L -o "$GO_TAR" "$GO_URL"

echo "解压到 ${INSTALL_DIR}..."
mkdir -p "$INSTALL_DIR"
tar -C "$INSTALL_DIR" -xzf "$GO_TAR"

echo "验证..."
"$INSTALL_DIR/go/bin/go" version

echo "完成。Go 已安装到: $INSTALL_DIR/go"
echo "Cursor 已配置为使用: $INSTALL_DIR/go/bin/go 安装工具。"
echo "请在 Cursor 中执行: Ctrl+Shift+P -> Go: Install/Update Tools -> 勾选 gopls -> OK"
