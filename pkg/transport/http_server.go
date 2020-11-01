package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewHTTPServer creates a preconfigured http.Server, eventually routes must be injected to the instance
func NewHTTPServer(lc fx.Lifecycle, r *mux.Router, logger *zap.Logger, cfg configuration.Configuration) *http.Server {
	addr := getHTTPAddress(cfg)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      &ochttp.Handler{Handler: r},
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				logger.Info("started http server", zap.String("address", addr),
					zap.String("endpoint", cfg.HTTP.Endpoint))
				_ = server.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopped http server", zap.String("address", addr),
				zap.String("endpoint", cfg.HTTP.Endpoint))
			return server.Shutdown(ctx)
		},
	})

	return server
}
