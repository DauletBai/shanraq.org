// github.com/DauletBai/shanraq.org/http/server.go
package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	k "github.com/DauletBai/shanraq.org/core/kernel"
)

// Server wraps the standard http.Server and integrates with Shanraq components.
type Server struct {
	kernel     *k.Kernel
	httpServer *http.Server
	router     *Router 
}

// NewServer creates a new Shanraq HTTP Server.
func NewServer(kernel *k.Kernel, router *Router) *Server {
	addr := kernel.Config().GetString("server.address") // e.g., "localhost" or "" for all interfaces
	port := kernel.Config().GetInt("server.port")
	if port == 0 {
		port = 8080 // Default port if not specified
	}
	fullAddr := fmt.Sprintf("%s:%d", addr, port)

	readTimeoutSeconds := kernel.Config().GetInt("server.read_timeout_seconds")
	if readTimeoutSeconds == 0 {
		readTimeoutSeconds = 5 // Default
	}
	writeTimeoutSeconds := kernel.Config().GetInt("server.write_timeout_seconds")
	if writeTimeoutSeconds == 0 {
		writeTimeoutSeconds = 10 // Default
	}
	idleTimeoutSeconds := kernel.Config().GetInt("server.idle_timeout_seconds")
	if idleTimeoutSeconds == 0 {
		idleTimeoutSeconds = 120 // Default
	}

	return &Server{
		kernel: kernel,
		router: router,
		httpServer: &http.Server{
			Addr:         fullAddr,
			Handler:      router, 
			ReadTimeout:  time.Duration(readTimeoutSeconds) * time.Second,
			WriteTimeout: time.Duration(writeTimeoutSeconds) * time.Second,
			IdleTimeout:  time.Duration(idleTimeoutSeconds) * time.Second,
		},
	}
}

// ListenAndServe starts the HTTP server and handles graceful shutdown.
func (s *Server) ListenAndServe() error {
	logger := s.kernel.Logger()
	logger.Info("Starting HTTP server...", "address", s.httpServer.Addr)

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info(fmt.Sprintf("HTTP server listening on %s", s.httpServer.Addr))
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "error", err)
			return err
		}
		logger.Info("HTTP server shut down.") // Log if ErrServerClosed or no error
		return nil                            // nil if closed gracefully or no error

	case sig := <-quit:
		logger.Info("OS signal received, initiating graceful shutdown...", "signal", sig.String())

		shutdownTimeout := s.kernel.Config().GetInt("server.shutdown_timeout_seconds")
		if shutdownTimeout == 0 {
			shutdownTimeout = 15 // Default
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			logger.Error("HTTP server graceful shutdown failed", "error", err)
			return err
		}
		logger.Info("HTTP server gracefully shut down.")
	}

	return nil
}