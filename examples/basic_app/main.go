// github.com/DauletBai/shanraq.org/examples/basic_app/main.go
package main

import (
	//"database/sql"
	"log"      // Standard log for initial error checking before framework logger is ready
	"net/http" // Для http.StatusOK
	"os"
	"strings"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware" 

	"github.com/DauletBai/shanraq.org/core/config"
	"github.com/DauletBai/shanraq.org/core/kernel"
	"github.com/DauletBai/shanraq.org/core/logger"
	shqHTTP "github.com/DauletBai/shanraq.org/http"
	frameworkMiddleware "github.com/DauletBai/shanraq.org/http/middleware"
)

func main() {
	var _ = config.DefaultConfigFileName
	var _ logger.LogLevel = logger.LevelDebug

	dummyConfigContent := `
app:
  name: "My Shanraq App"
  environment: "development" 
  version: "1.0.0"

server:
  address: "localhost"
  port: 8080
  read_timeout_seconds: 10
  write_timeout_seconds: 10
  idle_timeout_seconds: 60
  shutdown_timeout_seconds: 15 

logger:
  level: "DEBUG" 
  json_output: false
  add_source: true 

database:
  driver: "postgres" # or "pgx" if we use pgx directly without database/sql
  dsn: "postgres://youruser:yourpassword@localhost:5432/yourdbname?sslmode=disable"

  # Or you can use individual parameters:
  # host: "localhost"
  # port: 5432
  # user: "youruser"
  # password: "yourpassword"
  # dbname: "yourdbname"
  # sslmode: "disable" # или "require", "verify-full", etc.
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime_minutes: 2
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

	appRouter.GET("/test_validation", func(c *shqHTTP.Context) {
		username := c.QueryParam("username") 

		validationErrs := shqHTTP.NewValidationErrors() 

		if strings.TrimSpace(username) == "" {
			validationErrs.Add("username", "Username is required and cannot be empty.")
		} else if len(username) < 3 {
			validationErrs.Add("username", "Username must be at least 3 characters long.")
		}
		// You can add other checks, for example, for email, password, etc.
		// email := c.QueryParam("email")
		// if strings.TrimSpace(email) == "" {
		//     validationErrs.Add("email", "Email is required.")
		// }

		if !validationErrs.IsEmpty() {
			c.Logger().Info("Validation failed for /test_validation", "errors", validationErrs)
			// Передаем validationErrs в Details нашего HTTPError
			// Используем http.StatusBadRequest (400) или http.StatusUnprocessableEntity (422) для ошибок валидации
			err := shqHTTP.NewHTTPError(http.StatusBadRequest, "Input validation failed.", validationErrs)
			c.Error(err) // Наш Context.Error() отправит это как JSON
			return
		}

		// Если валидация прошла успешно
		c.JSON(http.StatusOK, map[string]string{
			"message":  "Validation successful!",
			"username": username,
		})
	})

	appRouter.GET("/db_test", func(c *shqHTTP.Context) {
		db := c.Kernel().DB() // Получаем доступ к *sql.DB из ядра
		if db == nil {
			c.Error(shqHTTP.NewHTTPError(http.StatusInternalServerError, "Database connection is not available."))
			return
		}

		// Простой запрос для получения текущего времени из БД
		var currentTime time.Time
		err := db.QueryRow("SELECT NOW()").Scan(&currentTime)
		if err != nil {
			c.Logger().Error("Database query failed", "error", err)
			c.Error(shqHTTP.NewHTTPError(http.StatusInternalServerError, "Failed to query database."))
			return
		}

		// Пример запроса к таблице (предположим, у вас есть таблица 'users' с полем 'username')
		// Этот код закомментирован, так как у нас нет такой таблицы в общем примере
		/*
		var username sql.NullString
		userId := 1 // Пример ID
		err = db.QueryRow("SELECT username FROM users WHERE id = $1", userId).Scan(&username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Logger().Info("User not found for db_test", "user_id", userId)
				c.Error(shqHTTP.NewHTTPError(http.StatusNotFound, "User not found."))
				return
			}
			c.Logger().Error("Database query for user failed", "error", err, "user_id", userId)
			c.Error(shqHTTP.NewHTTPError(http.StatusInternalServerError, "Failed to query user data."))
			return
		}
		*/
		
		c.JSON(http.StatusOK, map[string]interface{}{
			"message":      "Database test successful!",
			"current_time_from_db": currentTime.Format(time.RFC3339),
			// "found_username": username.String, // Раскомментируйте, если используете пример с users
		})
	})

    httpServer := shqHTTP.NewServer(appKernel, appRouter)
    appLogger.Info("Attempting to start HTTP server...")
    if err := httpServer.ListenAndServe(); err != nil {
        appLogger.Error("Failed to start or run HTTP server", "error", err)
        os.Exit(1)
    }
    appLogger.Info("Application has shut down gracefully.")
}