# Tenge Dev Tools Makefile
# Инструменты разработчика для языка Tenge

.PHONY: all install build test lint format clean help

# Default target
all: install

# Install all developer tools
install:
	@echo "🔧 Installing Tenge developer tools..."
	@echo "📦 Installing VS Code extension..."
	@if [ -d "vscode" ]; then \
		cd vscode && npm install; \
	fi
	@echo "📦 Installing Sublime Text package..."
	@if [ -d "sublime" ]; then \
		cp sublime/tenge.sublime-syntax ~/.config/sublime-text/Packages/User/; \
		cp sublime/tenge.sublime-snippets ~/.config/sublime-text/Packages/User/; \
	fi
	@echo "📦 Installing CLI tools..."
	@chmod +x linter/tenge-linter.py
	@chmod +x formatter/tenge-formatter.py
	@chmod +x cli/tenge-cli.py
	@echo "✅ Developer tools installed successfully!"

# Build VS Code extension
build-vscode:
	@echo "🔨 Building VS Code extension..."
	@if [ -d "vscode" ]; then \
		cd vscode && npm run compile; \
	fi

# Build all tools
build: build-vscode
	@echo "🔨 Building all developer tools..."

# Test all tools
test:
	@echo "🧪 Testing Tenge developer tools..."
	@echo "🔧 Testing linter..."
	@python3 linter/tenge-linter.py --help
	@echo "🔧 Testing formatter..."
	@python3 formatter/tenge-formatter.py --help
	@echo "🔧 Testing CLI..."
	@python3 cli/tenge-cli.py --help
	@echo "✅ All tools tested successfully!"

# Lint project files
lint:
	@echo "🔍 Linting Tenge project files..."
	@find . -name "*.tng" -exec python3 linter/tenge-linter.py {} \;

# Format project files
format:
	@echo "🎨 Formatting Tenge project files..."
	@find . -name "*.tng" -exec python3 formatter/tenge-formatter.py -i {} \;

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@if [ -d "vscode/out" ]; then \
		rm -rf vscode/out; \
	fi
	@if [ -d "vscode/node_modules" ]; then \
		rm -rf vscode/node_modules; \
	fi
	@echo "✅ Clean completed!"

# Create new Tenge project
new-project:
	@echo "🚀 Creating new Tenge project..."
	@read -p "Enter project name: " project_name; \
	python3 cli/tenge-cli.py init $$project_name

# Build existing Tenge project
build-project:
	@echo "🔨 Building Tenge project..."
	@python3 cli/tenge-cli.py build

# Test existing Tenge project
test-project:
	@echo "🧪 Testing Tenge project..."
	@python3 cli/tenge-cli.py test

# Run Tenge project
run-project:
	@echo "🚀 Running Tenge project..."
	@python3 cli/tenge-cli.py run

# Lint Tenge project
lint-project:
	@echo "🔍 Linting Tenge project..."
	@python3 cli/tenge-cli.py lint

# Format Tenge project
format-project:
	@echo "🎨 Formatting Tenge project..."
	@python3 cli/tenge-cli.py format

# Check dependencies
check-deps:
	@echo "🔍 Checking dependencies..."
	@./check_dependencies.sh

# Install VS Code extension
install-vscode: check-deps
	@echo "📦 Installing VS Code extension..."
	@if [ -d "vscode" ]; then \
		cd vscode && npm install && npm run compile && npm run package; \
		code --install-extension tenge-language-support-1.0.0.vsix; \
	fi

# Quick VS Code installation
install-vscode-quick: check-deps
	@echo "🚀 Quick VS Code installation..."
	@./install_vscode.sh

# Install Sublime Text package
install-sublime:
	@echo "📦 Installing Sublime Text package..."
	@if [ -d "sublime" ]; then \
		mkdir -p ~/.config/sublime-text/Packages/User/; \
		cp sublime/tenge.sublime-syntax ~/.config/sublime-text/Packages/User/; \
		cp sublime/tenge.sublime-snippets ~/.config/sublime-text/Packages/User/; \
	fi

# Create development environment
dev-env:
	@echo "🛠️  Setting up development environment..."
	@echo "📦 Installing Python dependencies..."
	@pip3 install -r requirements.txt
	@echo "📦 Installing Node.js dependencies..."
	@if [ -d "vscode" ]; then \
		cd vscode && npm install; \
	fi
	@echo "✅ Development environment ready!"

# Show help
help:
	@echo "Tenge Developer Tools"
	@echo "===================="
	@echo ""
	@echo "Available commands:"
	@echo "  install          - Install all developer tools"
	@echo "  build            - Build all tools"
	@echo "  test             - Test all tools"
	@echo "  lint             - Lint project files"
	@echo "  format           - Format project files"
	@echo "  clean            - Clean build artifacts"
	@echo "  new-project      - Create new Tenge project"
	@echo "  build-project    - Build existing Tenge project"
	@echo "  test-project     - Test existing Tenge project"
	@echo "  run-project      - Run Tenge project"
	@echo "  lint-project     - Lint Tenge project"
	@echo "  format-project   - Format Tenge project"
	@echo "  check-deps           - Check dependencies (Node.js, VS Code)"
	@echo "  install-vscode        - Install VS Code extension"
	@echo "  install-vscode-quick  - Quick VS Code installation"
	@echo "  install-sublime       - Install Sublime Text package"
	@echo "  dev-env          - Set up development environment"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make install           # Install all tools"
	@echo "  make new-project       # Create new project"
	@echo "  make build-project     # Build project"
	@echo "  make test-project      # Test project"
	@echo "  make run-project       # Run project"
