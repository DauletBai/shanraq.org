# Makefile для запуска бенчмарков Shanraq.org
# Shanraq.org бенчмарктарды іске қосу үшін Makefile

.PHONY: all benchmarks advanced clean results help

# Запуск всех бенчмарков
all: benchmarks

# Запуск базовых бенчмарков
benchmarks:
	@echo "🚀 Запуск Shanraq.org базовых бенчмарков..."
	@echo "📊 Генерация SVG результатов..."
	@./generate_demo_svgs.sh
	@echo "✅ Базовые бенчмарки завершены!"

# Запуск продвинутых бенчмарков
advanced:
	@echo "🚀 Запуск Shanraq.org продвинутых бенчмарков..."
	@echo "📊 Генерация продвинутых SVG результатов..."
	@echo "🌐 Сетевые оптимизации: epoll/kqueue + edge-triggered + ring-буферы..."
	@echo "📄 SIMD JSON: Stage-1/Stage-2 pipeline + runtime dispatch..."
	@echo "🔢 Матричные операции: CPU тайлинг + GPU shared-memory..."
	@echo "🧵 Concurrency: Lock-free + work-stealing + tail-latency..."
	@echo "📊 Генерация продвинутых SVG файлов..."
	@./generate_advanced_svgs.sh
	@echo "✅ Продвинутые бенчмарки завершены!"

# Очистка результатов
clean:
	@echo "🧹 Очистка результатов..."
	@rm -rf results/
	@rm -f *.svg
	@rm -f COMPREHENSIVE_BENCHMARK_REPORT.md
	@rm -f ADVANCED_COMPREHENSIVE_REPORT.md
	@echo "✅ Очистка завершена"

# Показать результаты
results:
	@echo "📊 Результаты бенчмарков:"
	@ls -la results/ 2>/dev/null || echo "  Результаты не найдены. Запустите 'make benchmarks' или 'make advanced'"

# Помощь
help:
	@echo "Shanraq.org Advanced Benchmark Suite"
	@echo "===================================="
	@echo "Доступные команды:"
	@echo "  make benchmarks  - Запуск базовых бенчмарков"
	@echo "  make advanced    - Запуск продвинутых бенчмарков"
	@echo "  make clean       - Очистка результатов"
	@echo "  make results     - Показать результаты"
	@echo "  make help        - Показать эту справку"
	@echo ""
	@echo "Продвинутые бенчмарки включают:"
	@echo "  🌐 Сетевые оптимизации: epoll/kqueue + edge-triggered + ring-буферы"
	@echo "  📄 SIMD JSON: Stage-1/Stage-2 pipeline + runtime dispatch"
	@echo "  🔢 Матричные операции: CPU тайлинг + GPU shared-memory"
	@echo "  🧵 Concurrency: Lock-free + work-stealing + tail-latency"