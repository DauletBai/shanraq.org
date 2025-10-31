package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"shanraq.org/pkg/shanraq"
)

//go:embed sql/*.sql
var embedded embed.FS

// Module runs embedded goose migrations against the shared DB pool.
type Module struct {
	fs  fs.FS
	dir string
}

// New returns a migration module backed by the internal migration set.
func New() *Module {
	return &Module{
		fs:  embedded,
		dir: "sql",
	}
}

func (m *Module) Name() string {
	return "migrations"
}

// Init executes migrations before any HTTP traffic or workers start.
func (m *Module) Init(ctx context.Context, rt *shanraq.Runtime) error {
	goose.SetLogger(goose.NopLogger()) // silence CLI noise in app runtime
	goose.SetBaseFS(m.fs)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := ensureDatabase(ctx, rt.Config.Database.URL); err != nil {
		return fmt.Errorf("ensure database: %w", err)
	}

	sqlDB, err := sql.Open("pgx", rt.Config.Database.URL)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	defer sqlDB.Close()

	if err := goose.UpContext(ctx, sqlDB, m.dir); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	rt.Logger.Info("migrations applied")
	return nil
}

var _ interface {
	shanraq.Module
	shanraq.InitializerModule
} = (*Module)(nil)

func ensureDatabase(parentCtx context.Context, dsn string) error {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse database url: %w", err)
	}

	targetDB := cfg.ConnConfig.Database
	if targetDB == "" {
		return nil
	}

	if err := waitForDatabase(parentCtx, cfg.ConnConfig, 8); err == nil {
		return nil
	} else if !isMissingDatabase(err) {
		return err
	}

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	adminConf := cfg.ConnConfig.Copy()
	if adminConf.Database == "" {
		adminConf.Database = "postgres"
	} else {
		adminConf.Database = "postgres"
	}

	adminConn, err := pgx.ConnectConfig(ctx, adminConf)
	if err != nil {
		return fmt.Errorf("connect postgres database: %w", err)
	}
	defer adminConn.Close(ctx)

	createStmt := fmt.Sprintf("CREATE DATABASE %s", pgx.Identifier{targetDB}.Sanitize())
	if _, err := adminConn.Exec(ctx, createStmt); err != nil && !isDuplicateDatabase(err) {
		return fmt.Errorf("create database %s: %w", targetDB, err)
	}

	return nil
}

func pingDatabase(ctx context.Context, conf *pgx.ConnConfig) error {
	conn, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	return conn.Ping(ctx)
}

func isMissingDatabase(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "3D000"
	}
	return false
}

func isDuplicateDatabase(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P04"
	}
	return false
}

func waitForDatabase(parentCtx context.Context, conf *pgx.ConnConfig, attempts int) error {
	var err error
	backoff := time.Second
	for attempt := 1; attempt <= attempts; attempt++ {
		ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
		err = pingDatabase(ctx, conf)
		cancel()
		if err == nil || isMissingDatabase(err) {
			return err
		}
		time.Sleep(backoff)
		if backoff < 5*time.Second {
			backoff *= 2
		}
	}
	return err
}
