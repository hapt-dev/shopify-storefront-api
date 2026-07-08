// Package server builds and runs the HTTP server on top of the infrastructure from package app.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/app"
)

const (
	shutdownTimeout   = 15 * time.Second
	readHeaderTimeout = 10 * time.Second
)

// Server wraps http.Server together with the shared infrastructure (app) to manage its lifecycle.
type Server struct {
	httpServer *http.Server
	app        *app.App
}

// New takes the already-bootstrapped infrastructure (app) and builds the controllers, router and http.Server.
func New(a *app.App) *Server {
	controllers := newControllers(a.Services)
	router := NewRouter(a.Log, controllers)

	return &Server{
		httpServer: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", a.Config.Host, a.Config.Port),
			Handler:           router,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		app: a,
	}
}

// Run starts the server and blocks until SIGINT/SIGTERM is received,
// then gracefully shuts down: stops the http server and closes the infrastructure (app.Close).
func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := s.app.Log

	serverErr := make(chan error, 1)
	go func() {
		log.Info("http server starting", zap.String("addr", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("http server shutdown error", zap.Error(err))
	}

	s.app.Close(shutdownCtx)
	log.Info("http server stopped gracefully")
	return nil
}
