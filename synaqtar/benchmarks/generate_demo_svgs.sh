#!/bin/bash
# generate_demo_svgs.sh - Генерация демонстрационных SVG файлов
# Демонстрациялық SVG файлдарын генерациялау

echo "🎨 Генерация демонстрационных SVG файлов..."
echo "=========================================="

# Создание директории результатов
mkdir -p results

# Получение текущего времени
TIMESTAMP=$(date +"%Y.%m.%d_%H:%M")

echo "⏰ Время: $TIMESTAMP"
echo "📁 Директория results создана"

# Генерация SVG файлов
echo ""
echo "📊 Генерация SVG файлов..."

# 1. Monte Carlo
cat > "results/Monte_Carlo_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">SIMD Monte Carlo Pi Estimation</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 1250.5 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Iterations: 10,000,000</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 8,000,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 45.2 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#27ae60">SIMD Acceleration: 4.2x</text>
<text x="70" y="280" font-family="Arial" font-size="14" fill="#34495e">Accuracy: 99.99%</text>
<rect x="300" y="200" width="100" height="200" fill="#3498db"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: Monte_Carlo_${TIMESTAMP}.svg"

# 2. Fibonacci
cat > "results/Fibonacci_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">SIMD Fibonacci Benchmark</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 85.3 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Input Size: 40</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 117,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 12.8 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#27ae60">SIMD Acceleration: 3.8x</text>
<rect x="300" y="250" width="100" height="150" fill="#e74c3c"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: Fibonacci_${TIMESTAMP}.svg"

# 3. QuickSort
cat > "results/QuickSort_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">SIMD QuickSort Benchmark</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 340.7 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Array Size: 100,000</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 293,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 25.6 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#27ae60">SIMD Acceleration: 2.9x</text>
<rect x="300" y="220" width="100" height="180" fill="#f39c12"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: QuickSort_${TIMESTAMP}.svg"

# 4. Matrix Multiplication
cat > "results/Matrix_Multiplication_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">SIMD Matrix Multiplication</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 1850.2 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Matrix Size: 500x500</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 270,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 78.4 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#27ae60">SIMD Acceleration: 6.1x</text>
<rect x="300" y="150" width="100" height="250" fill="#9b59b6"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: Matrix_Multiplication_${TIMESTAMP}.svg"

# 5. SIMD JSON Parsing
cat > "results/SIMD_JSON_Parsing_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">SIMD JSON Parsing</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 450.8 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">JSON Size: 1,000,000 bytes</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 2,220</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 15.3 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#27ae60">SIMD Acceleration: 3.5x</text>
<rect x="300" y="200" width="100" height="200" fill="#1abc9c"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: SIMD_JSON_Parsing_${TIMESTAMP}.svg"

# 6. Zero-Copy HTTP Requests
cat > "results/Zero_Copy_HTTP_Requests_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">Zero-Copy HTTP Requests</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 450.8 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Requests: 10,000</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 22,200</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 15.3 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#1abc9c">Zero-Copy Efficiency: 95.7%</text>
<text x="70" y="280" font-family="Arial" font-size="14" fill="#34495e">Throughput: 125.6 MB/s</text>
<text x="70" y="300" font-family="Arial" font-size="14" fill="#34495e">Latency: 0.045 ms</text>
<rect x="300" y="200" width="100" height="200" fill="#e67e22"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: Zero_Copy_HTTP_Requests_${TIMESTAMP}.svg"

# 7. GPU Matrix Multiplication
cat > "results/GPU_Matrix_Multiplication_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">GPU Matrix Multiplication</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 125.4 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Matrix Size: 1000x1000</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 8,000,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 78.4 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#8e44ad">GPU Acceleration: 15.2x</text>
<rect x="300" y="300" width="100" height="100" fill="#8e44ad"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: GPU_Matrix_Multiplication_${TIMESTAMP}.svg"

# 8. TLS Benchmark
cat > "results/TLS_Benchmark_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">Thread-Local Storage Benchmark</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 250.6 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Iterations: 100,000</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 399,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 8.2 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#e67e22">Threading Efficiency: 87.5%</text>
<rect x="300" y="250" width="100" height="150" fill="#e67e22"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: TLS_Benchmark_${TIMESTAMP}.svg"

# 9. Thread Pool Benchmark
cat > "results/Thread_Pool_Benchmark_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">Thread Pool Benchmark</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 180.3 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Tasks: 1,000,000</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 5,550,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 12.4 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#e67e22">Threading Efficiency: 92.3%</text>
<rect x="300" y="280" width="100" height="120" fill="#e67e22"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: Thread_Pool_Benchmark_${TIMESTAMP}.svg"

# 10. Message Passing Benchmark
cat > "results/Message_Passing_Benchmark_${TIMESTAMP}.svg" << 'EOF'
<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
<rect width="800" height="600" fill="#f8f9fa"/>
<text x="400" y="50" text-anchor="middle" font-family="Arial" font-size="24" font-weight="bold" fill="#2c3e50">Message Passing Benchmark</text>
<text x="400" y="80" text-anchor="middle" font-family="Arial" font-size="14" fill="#7f8c8d">Shanraq.org Runtime - 2025.01.11_15:30</text>
<rect x="50" y="120" width="700" height="400" fill="white" stroke="#bdc3c7" stroke-width="2"/>
<text x="70" y="150" font-family="Arial" font-size="16" font-weight="bold" fill="#2c3e50">Performance Results:</text>
<text x="70" y="180" font-family="Arial" font-size="14" fill="#34495e">Execution Time: 320.7 ms</text>
<text x="70" y="200" font-family="Arial" font-size="14" fill="#34495e">Messages: 100,000</text>
<text x="70" y="220" font-family="Arial" font-size="14" fill="#34495e">Operations/sec: 312,000</text>
<text x="70" y="240" font-family="Arial" font-size="14" fill="#34495e">Memory Usage: 6.8 MB</text>
<text x="70" y="260" font-family="Arial" font-size="14" fill="#e67e22">Threading Efficiency: 89.1%</text>
<rect x="300" y="230" width="100" height="170" fill="#e67e22"/>
<text x="350" y="420" text-anchor="middle" font-family="Arial" font-size="12" fill="#2c3e50">Performance Bar</text>
</svg>
EOF
echo "  ✅ Создан: Message_Passing_Benchmark_${TIMESTAMP}.svg"

echo ""
echo "🎉 Генерация завершена!"
echo "📁 SVG файлы сохранены в: synaqtar/benchmarks/results/"
echo "📊 Всего создано: 10 SVG файлов"


