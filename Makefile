# Shanraq.org Makefile
# Agglutinative Web Application Management

.PHONY: help build run test clean install demo benchmark benchmark-financial benchmark-crud benchmark-network benchmark-comparison benchmark-verbose show-benchmark-results clean-benchmarks

# Default target
help:
	@echo "🚀 Shanraq.org - Agglutinative Web Application"
	@echo "=========================================="
	@echo ""
	@echo "Available commands:"
	@echo "  make help           - Show this help message"
	@echo "  make build          - Build the project"
	@echo "  make run            - Run the server"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-e2e       - Run E2E tests"
	@echo "  make lint           - Check code for errors"
	@echo "  make format         - Format code"
	@echo "  make install        - Install dependencies"
	@echo "  make install-optimized - Install optimized dependencies"
	@echo "  make clean          - Clean the project"
	@echo "  make demo           - Run demo server"
	@echo "  make benchmark      - Run full benchmark"
	@echo "  make benchmark-financial - Run financial-mathematical benchmarks"
	@echo "  make benchmark-crud - Run CRUD database benchmarks"
	@echo "  make benchmark-network - Run network benchmarks"
	@echo "  make benchmark-comparison - Run language comparison benchmarks"
	@echo "  make benchmark-verbose - Run detailed benchmark"
	@echo "  make show-benchmark-results - Show benchmark results"
	@echo "  make clean-benchmarks - Clean benchmark results"
	@echo ""

build:
	@echo "🔨 Building Shanraq.org..."
	@npm run build

run:
	@echo "🚀 Starting Shanraq.org..."
	@npm start

test: test-unit test-integration test-e2e

test-unit:
	@echo "🧪 Running unit tests..."
	@npm run test:unit

test-integration:
	@echo "🧪 Running integration tests..."
	@npm run test:integration

test-e2e:
	@echo "🧪 Running E2E tests..."
	@npm run test:e2e

lint:
	@echo "🔍 Checking code..."
	@npm run lint

format:
	@echo "🎨 Formatting code..."
	@npm run format

install:
	@echo "📦 Installing Shanraq.org..."
	@npm install

install-optimized:
	@echo "📦 Installing optimized dependencies..."
	@cp package_optimized.json package.json
	@npm install
	@echo "✅ Optimized dependencies installed!"

clean:
	@echo "🧹 Cleaning project..."
	@rm -rf node_modules build dist coverage

demo:
	@echo "🎯 Starting Shanraq.org demo server..."
	@echo "=========================================="
	@echo "🌐 Server: http://localhost:8080"
	@echo "📄 Home: http://localhost:8080/"
	@echo "📝 Blog: http://localhost:8080/blog"
	@echo "👥 About: http://localhost:8080/about"
	@echo "📞 Contact: http://localhost:8080/contact"
	@echo "🔧 API: http://localhost:8080/api/v1/health"
	@echo "=========================================="
	@echo "Press Ctrl+C to stop the server"
	@echo ""
	python3 synaqtar/demo/ulgi_server.py

benchmark:
	@echo "🚀 Running Shanraq.org benchmark tests..."
	@echo "====================================================="
	@cd synaqtar/benchmarks && make benchmark-full

benchmark-financial:
	@echo "📊 Financial-mathematical benchmarks..."
	@cd synaqtar/benchmarks && make benchmark-financial

benchmark-crud:
	@echo "💾 CRUD database benchmarks..."
	@cd synaqtar/benchmarks && make benchmark-crud

benchmark-network:
	@echo "🌐 Network benchmarks..."
	@cd synaqtar/benchmarks && make benchmark-network

benchmark-comparison:
	@echo "🔍 Language comparison benchmarks..."
	@cd synaqtar/benchmarks && make benchmark-comparison

benchmark-verbose:
	@echo "🚀 Detailed Shanraq.org benchmark..."
	@cd synaqtar/benchmarks && make benchmark-verbose

show-benchmark-results:
	@echo "📊 Shanraq.org benchmark results"
	@cd synaqtar/benchmarks && make show-results

clean-benchmarks:
	@echo "🧹 Cleaning benchmark results..."
	@cd synaqtar/benchmarks && make clean-benchmarks