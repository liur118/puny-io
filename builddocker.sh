#!/bin/bash

set -e

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# 镜像名称
IMAGE_NAME="puny-io:latest"

# 构建 docker 镜像
cd "$SCRIPT_DIR"
echo "正在构建 Docker 镜像 $IMAGE_NAME ..."
docker build -t $IMAGE_NAME .
echo "Docker 镜像构建完成: $IMAGE_NAME" 