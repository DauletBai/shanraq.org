// File: shanraq.org/database/db.go
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	//"github.com/DauletBai/shanraq.org/core/config"
	//"github.com/DauletBai/shanraq.org/core/logger"
)

// ConfigProvider определяет методы, необходимые для получения настроек БД.
// Это подмножество вашего существующего config.Provider.
type ConfigProvider interface {
	GetString(key string) string
	GetInt(key string) int
	IsSet(key string) bool 
}

// LoggerInterface определяет методы, необходимые для логирования.
// Это ваш существующий logger.Interface.
type LoggerInterface interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{}) 
	Error(msg string, args ...interface{})
}


// NewConnection теперь принимает ConfigProvider и LoggerInterface вместо всего Kernel.
func NewConnection(cfg ConfigProvider, logger LoggerInterface) (*sql.DB, error) {
	driver := cfg.GetString("database.driver")
	dsn := cfg.GetString("database.dsn")

	if dsn == "" {
		host := cfg.GetString("database.host")
		port := cfg.GetInt("database.port")
		user := cfg.GetString("database.user")
		password := cfg.GetString("database.password") 
		dbname := cfg.GetString("database.dbname")
		sslmode := cfg.GetString("database.sslmode")

		if host == "" || port == 0 || user == "" || dbname == "" {
			return nil, fmt.Errorf("database DSN is not configured and individual parameters (host, port, user, dbname) are missing or incomplete")
		}
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	}

	if driver == "" {
		driver = "pgx"
		logger.Info("Database driver not specified, defaulting to 'pgx' (for database/sql).")
	} else if driver == "postgres" {
		driver = "pgx"
	}

	logger.Info("Attempting to connect to database...", "driver", driver)
	// logger.Debug("Database DSN", "dsn", dsn) // Осторожно с логированием DSN

	db, err := sql.Open(driver, dsn)
	if err != nil {
		logger.Error("Failed to open database connection", "error", err)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	maxOpenConns := cfg.GetInt("database.max_open_conns")
	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
	} else {
		db.SetMaxOpenConns(25)
	}

	maxIdleConns := cfg.GetInt("database.max_idle_conns")
	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns)
	} else {
		db.SetMaxIdleConns(25)
	}

	connMaxLifetimeMinutes := cfg.GetInt("database.conn_max_lifetime_minutes")
	if connMaxLifetimeMinutes > 0 {
		db.SetConnMaxLifetime(time.Duration(connMaxLifetimeMinutes) * time.Minute)
	} else {
		db.SetConnMaxLifetime(5 * time.Minute)
	}

	if err = db.Ping(); err != nil {
		logger.Error("Failed to ping database", "error", err)
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to the database.")
	return db, nil
}