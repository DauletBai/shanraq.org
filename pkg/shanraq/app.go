package shanraq

import (
	"context"
	"fmt"

	"shanraq.org/internal/config"
	"shanraq.org/internal/db"
	"shanraq.org/internal/httpserver"
	"shanraq.org/internal/logging"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Module provides a named building block that can hook into runtime stages.
type Module interface {
	Name() string
}

// RouterModule registers HTTP routes onto the shared router.
type RouterModule interface {
	Module
	Routes(r chi.Router)
}

// StarterModule is invoked once the runtime is ready; it can block or spawn workers.
type StarterModule interface {
	Module
	Start(ctx context.Context, rt *Runtime) error
}

// InitializerModule runs before the HTTP server starts, ideal for validation or DB bootstrap.
type InitializerModule interface {
	Module
	Init(ctx context.Context, rt *Runtime) error
}

// Runtime is provided to modules with the fully bootstrapped dependencies.
type Runtime struct {
	Config config.Config
	Logger *zap.Logger
	DB     *pgxpool.Pool
	Router chi.Router
}

// Application wires together configuration, dependencies, and modules.
type Application struct {
	cfg     config.Config
	modules []Module
}

// New builds the application with the provided configuration.
func New(cfg config.Config) *Application {
	return &Application{cfg: cfg}
}

// Register attaches a module to the application.
func (a *Application) Register(mod Module) {
	a.modules = append(a.modules, mod)
}

// Run boots all dependencies and blocks until the context is canceled or an error occurs.
func (a *Application) Run(ctx context.Context) error {
	logger, err := logging.Build(a.cfg.Logging)
	if err != nil {
		return err
	}
	defer logger.Sync() //nolint:errcheck

	pool, err := db.Connect(ctx, a.cfg.Database, logger)
	if err != nil {
		return err
	}
	defer pool.Close()

	server := httpserver.New(a.cfg.Server, logger)
	rt := &Runtime{
		Config: a.cfg,
		Logger: logger,
		DB:     pool,
		Router: server.Router(),
	}

	for _, mod := range a.modules {
		if initializer, ok := mod.(InitializerModule); ok {
			if err := initializer.Init(ctx, rt); err != nil {
				return fmt.Errorf("%s init: %w", mod.Name(), err)
			}
		}

		if router, ok := mod.(RouterModule); ok {
			router.Routes(rt.Router)
		}
	}

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if err := server.Start(groupCtx); err != nil {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})

	for _, mod := range a.modules {
		module := mod
		if starter, ok := module.(StarterModule); ok {
			group.Go(func() error {
				if err := starter.Start(groupCtx, rt); err != nil {
					return fmt.Errorf("%s start: %w", module.Name(), err)
				}
				return nil
			})
		}
	}

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}
