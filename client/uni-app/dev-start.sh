#!/bin/bash

# 社区生活 uni-app 开发启动脚本

echo "🚀 社区生活 uni-app 开发环境启动"
echo "=================================="

# 检查 Node.js 版本
if ! command -v node &> /dev/null; then
    echo "❌ 错误：未安装 Node.js"
    echo "请访问 https://nodejs.org/ 下载安装"
    exit 1
fi

echo "✓ Node.js 版本: $(node -v)"
echo "✓ npm 版本: $(npm -v)"

# 检查依赖
if [ ! -d "node_modules" ]; then
    echo ""
    echo "📦 正在安装依赖..."
    npm install
    if [ $? -ne 0 ]; then
        echo "❌ 依赖安装失败"
        exit 1
    fi
    echo "✓ 依赖安装完成"
fi

echo ""
echo "请选择启动模式："
echo "1) H5 开发模式 (http://localhost:3000)"
echo "2) 微信小程序开发模式"
echo "3) 构建生产版本 - H5"
echo "4) 构建生产版本 - 微信小程序"
echo ""
read -p "请输入选项 (1-4): " choice

case $choice in
    1)
        echo ""
        echo "🌐 启动 H5 开发模式..."
        npm run dev:h5
        ;;
    2)
        echo ""
        echo "📱 启动微信小程序开发模式..."
        npm run dev:mp-weixin
        ;;
    3)
        echo ""
        echo "🔨 构建 H5 生产版本..."
        npm run build:h5
        ;;
    4)
        echo ""
        echo "🔨 构建微信小程序生产版本..."
        npm run build:mp-weixin
        ;;
    *)
        echo "❌ 无效选项"
        exit 1
        ;;
esac