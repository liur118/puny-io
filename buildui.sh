#!/bin/bash

# 设置错误时退出
set -e

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# 定义输出目录（相对于脚本目录）
BUILD_DIR="$SCRIPT_DIR/build/puny-io"
UI_DIR="$BUILD_DIR/ui"

# 进入 UI 目录（相对于脚本目录）
cd "$SCRIPT_DIR/ui"

# 安装依赖
echo "正在安装依赖..."
npm install

# 构建项目
echo "正在构建项目..."
npm run build

# 创建输出目录（如果不存在）
echo "正在创建输出目录..."
mkdir -p "$UI_DIR"

# 复制构建文件到输出目录
echo "正在复制构建文件..."
cp -r dist/* "$UI_DIR/"

# 返回脚本目录
cd "$SCRIPT_DIR"

echo "构建完成！UI 文件已输出到 $UI_DIR 目录" 