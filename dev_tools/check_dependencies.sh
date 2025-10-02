#!/bin/bash

# Tenge Dev Tools Dependency Checker
# Проверка зависимостей для Tenge Dev Tools

echo "🔍 Checking dependencies for Tenge Dev Tools..."

# Проверка Node.js
echo "📦 Checking Node.js..."
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo "✅ Node.js found: $NODE_VERSION"
else
    echo "❌ Node.js not found"
    echo "   Please install Node.js from: https://nodejs.org/"
    echo "   Or use Homebrew: brew install node"
    exit 1
fi

# Проверка npm
echo "📦 Checking npm..."
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm --version)
    echo "✅ npm found: $NPM_VERSION"
else
    echo "❌ npm not found"
    echo "   Please install npm (usually comes with Node.js)"
    exit 1
fi

# Проверка VS Code
echo "📦 Checking VS Code..."
if command -v code &> /dev/null; then
    CODE_VERSION=$(code --version | head -n1)
    echo "✅ VS Code found: $CODE_VERSION"
else
    echo "❌ VS Code not found"
    echo "   Please install VS Code from: https://code.visualstudio.com/"
    echo "   Or use Homebrew: brew install --cask visual-studio-code"
    echo ""
    echo "   After installation, add VS Code to PATH:"
    echo "   echo 'export PATH=\"/Applications/Visual Studio Code.app/Contents/Resources/app/bin:\$PATH\"' >> ~/.zshrc"
    echo "   source ~/.zshrc"
    exit 1
fi

echo ""
echo "🎉 All dependencies are installed!"
echo "🚀 You can now run: make install-vscode-quick"







