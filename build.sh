#!/bin/bash

# 设置错误时退出
set -e

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# 定义输出目录（相对于脚本目录）
BUILD_DIR="$SCRIPT_DIR/build/puny-io"
BIN_DIR="$BUILD_DIR/bin"
CONF_DIR="$BUILD_DIR/conf"

echo "开始构建项目..."

# 创建输出目录（如果不存在）
echo "正在创建输出目录..."
mkdir -p "$BIN_DIR"
mkdir -p "$CONF_DIR"

# 构建后端项目
echo "正在构建后端项目..."
go build -o "$BIN_DIR/puny-io" .

# 复制配置文件
echo "正在复制配置文件..."
cp -r conf/* "$CONF_DIR/"

# 构建前端项目
echo "正在构建前端项目..."
"$SCRIPT_DIR/buildui.sh"

# 返回脚本目录
cd "$SCRIPT_DIR"

echo "构建完成！"
echo "后端可执行文件位置: $BIN_DIR/puny-io"
echo "配置文件位置: $CONF_DIR"
echo "前端文件位置: $BUILD_DIR/ui"

# 打包整个 build/puny-io 目录为 tar.gz
TAR_PATH="$SCRIPT_DIR/build/puny-io.tar.gz"
echo "正在打包 $BUILD_DIR 到 $TAR_PATH ..."
tar -czf "$TAR_PATH" -C "$SCRIPT_DIR/build" puny-io

echo "打包完成！压缩包位置: $TAR_PATH"
