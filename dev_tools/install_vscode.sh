#!/bin/bash

# Tenge VS Code Extension Installer
# Установщик расширения VS Code для Tenge

set -e

echo "🚀 Installing Tenge VS Code Extension..."

# Проверка наличия VS Code
if ! command -v code &> /dev/null; then
    echo "❌ VS Code not found. Please install VS Code first."
    echo "   Download from: https://code.visualstudio.com/"
    exit 1
fi

# Проверка наличия Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js not found. Please install Node.js first."
    echo "   Download from: https://nodejs.org/"
    exit 1
fi

# Проверка наличия npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm not found. Please install npm first."
    exit 1
fi

echo "✅ VS Code and Node.js found"

# Переход в директорию расширения
cd "$(dirname "$0")/vscode"

echo "📦 Installing dependencies..."
npm install

echo "🔨 Compiling extension..."
npm run compile

echo "📦 Creating extension package..."
npm run package

echo "🔧 Installing extension in VS Code..."
code --install-extension tenge-language-support-1.0.0.vsix

echo "✅ Tenge VS Code Extension installed successfully!"
echo ""
echo "🎯 Next steps:"
echo "1. Restart VS Code"
echo "2. Open a .tng file"
echo "3. Enjoy Tenge language support!"
echo ""
echo "📚 For more information, see VSCODE_INSTALLATION.md"







