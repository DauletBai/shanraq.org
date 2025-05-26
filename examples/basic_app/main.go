// github.com/DauletBai/shanraq.org/examples/basic_app/main.go
package main

import (
	"log"      // Standard log for initial error checking before framework logger is ready
	"net/http" // Для http.StatusOK
	"os"

	chiMiddleware "github.com/go-chi/chi/v5/middleware" 

	"github.com/DauletBai/shanraq.org/core/config"
	"github.com/DauletBai/shanraq.org/core/kernel"
	"github.com/DauletBai/shanraq.org/core/logger"
	shqHTTP "github.com/DauletBai/shanraq.org/http"
	frameworkMiddleware "github.com/DauletBai/shanraq.org/http/middleware"
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

	appKernel, err := kernel.New (
		kernel.WithConfigFile("config", "./configs"), 
	)
	if err != nil {
		log.Fatalf("Failed to initialize Shanraq Kernel: %v", err)
	}

	appLogger := appKernel.Logger()
	appConfig := appKernel.Config() 

	appLogger.Info("Application starting...", "appName", appConfig.GetString("app.name"), "version", appConfig.GetString("app.version"))

    appRouter := shqHTTP.NewRouter(appKernel)

    appRouter.Use(chiMiddleware.RequestID)
    appRouter.Use(chiMiddleware.RealIP)
    appRouter.Use(chiMiddleware.Recoverer)

    appRouter.Use(frameworkMiddleware.RequestLogger(appKernel))

	appRouter.GET("/hello", func(c *shqHTTP.Context) { 
        c.Logger().Info("Handler for /hello called")
        c.JSON(http.StatusOK, map[string]string{
            "message": "Hello from Shanraq Framework!",
            "appName": c.Kernel().Config().GetString("app.name"),
        })
    })

	appRouter.GET("/panic", func(c *shqHTTP.Context) {
        c.Logger().Warn("This handler will panic!")
        panic("Simulated panic in handler")
    })

    httpServer := shqHTTP.NewServer(appKernel, appRouter)
    appLogger.Info("Attempting to start HTTP server...")
    if err := httpServer.ListenAndServe(); err != nil {
        appLogger.Error("Failed to start or run HTTP server", "error", err)
        os.Exit(1)
    }
    appLogger.Info("Application has shut down gracefully.")
}