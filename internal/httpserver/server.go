package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"shanraq.org/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Server wraps chi.Router with lifecycle hooks.
type Server struct {
	cfg    config.ServerConfig
	router chi.Router
	http   *http.Server
	logger *zap.Logger
}

// New instantiates the HTTP server with default middlewares wired up.
func New(cfg config.ServerConfig, logger *zap.Logger) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return &Server{
		cfg:    cfg,
		router: r,
		logger: logger,
	}
}

func (s *Server) Router() chi.Router {
	return s.router
}

// Start begins serving HTTP traffic until the context is canceled.
func (s *Server) Start(ctx context.Context) error {
	s.http = &http.Server{
		Addr:         s.cfg.Address,
		Handler:      s.router,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  s.cfg.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("http server starting", zap.String("addr", s.cfg.Address))
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("listen: %w", err)
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.http.Shutdown(shutdownCtx)
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.http == nil {
		return nil
	}
	return s.http.Shutdown(ctx)
}
