// github.com/DauletBai/shanraq.org/examples/basic_app/main.go
package main

import (
	//"fmt"
	"log"      // Standard log for initial error checking before framework logger is ready
	"net/http" // Для http.StatusOK
	"os"

	"github.com/go-chi/chi/v5/middleware" 

	"github.com/DauletBai/shanraq.org/core/config"
	"github.com/DauletBai/shanraq.org/core/kernel"
	"github.com/DauletBai/shanraq.org/core/logger"
	shanraqHTTP "github.com/DauletBai/shanraq.org/http"
)

func main() {
	// +++ Диагностические строки (оставляем, раз они помогают) +++
	var _ = config.DefaultConfigFileName
	var _ logger.LogLevel = logger.LevelDebug
	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	// --- Создание dummy config файла для примера ---
	// (Этот код остается таким же, как и раньше)
	dummyConfigContent := `
app:
  name: "My Shanraq App"
  environment: "development" 
  version: "1.0.0"

server:
  address: "localhost" # Адрес для прослушивания, пустая строка для всех интерфейсов
  port: 8080
  read_timeout_seconds: 10
  write_timeout_seconds: 10
  idle_timeout_seconds: 60
  shutdown_timeout_seconds: 15 # Таймаут для graceful shutdown

logger:
  level: "DEBUG" 
  json_output: false
  add_source: true # Включим исходники в логах для примера
`
	configDir := "./configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Failed to create config dir: %v", err)
	}
	configFilePath := configDir + "/config.yaml"
	if err := os.WriteFile(configFilePath, []byte(dummyConfigContent), 0644); err != nil {
		log.Fatalf("Failed to write dummy config: %v", err)
	}
	defer os.RemoveAll(configDir)

	// --- Инициализация ядра ---
	appKernel, err := kernel.New(
		kernel.WithConfigFile("config", "./configs"), // Загружаем наш dummy config
	)
	if err != nil {
		// Используем стандартный log, так как наш логгер может быть еще не инициализирован
		log.Fatalf("Failed to initialize Shanraq Kernel: %v", err)
	}

	appLogger := appKernel.Logger()
	appConfig := appKernel.Config() // Получаем доступ к конфигурации

	appLogger.Info("Application starting...", "appName", appConfig.GetString("app.name"), "version", appConfig.GetString("app.version"))

	// --- Настройка маршрутизатора ---
	appRouter := shanraqHTTP.NewRouter(appKernel)

	// Добавляем некоторые стандартные middleware от chi
	appRouter.Use(middleware.RequestID)  // Добавляет уникальный ID каждому запросу
	appRouter.Use(middleware.RealIP)     // Определяет реальный IP клиента (полезно за прокси)
	appRouter.Use(middleware.Recoverer)  // Перехватывает паники в обработчиках и возвращает 500

	// Простое middleware для логирования запросов с использованием нашего логгера
	appRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Создаем Shanraq контекст внутри middleware, если он нужен здесь.
			// В данном случае, мы используем логгер из ядра напрямую.
			reqLogger := appKernel.Logger().With("method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "requestID", middleware.GetReqID(r.Context()))
			reqLogger.Info("Incoming request")
			
			// Передаем управление следующему обработчику в цепочке
			next.ServeHTTP(w, r) 
		})
	})


	// Регистрируем тестовый маршрут
	appRouter.GET("/hello", func(c *shanraqHTTP.Context) {
		c.Logger().Info("Handler for /hello called") // Используем логгер из Shanraq Context
		c.JSON(http.StatusOK, map[string]string{
			"message": "Hello from Shanraq Framework!",
			"appName": c.Kernel().Config().GetString("app.name"), // Доступ к конфигу через контекст
		})
	})

	appRouter.GET("/panic", func(c *shanraqHTTP.Context) {
		c.Logger().Warn("This handler will panic!")
		panic("Simulated panic in handler")
	})

	// --- Создание и запуск HTTP сервера ---
	httpServer := shanraqHTTP.NewServer(appKernel, appRouter)

	appLogger.Info("Attempting to start HTTP server...")
	if err := httpServer.ListenAndServe(); err != nil {
		appLogger.Error("Failed to start or run HTTP server", "error", err)
		os.Exit(1) // Выходим, если сервер не смог запуститься
	}

	appLogger.Info("Application has shut down gracefully.")
}