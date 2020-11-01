package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// NewMux creates a preconfigured mux.Router, eventually routes must be injected to the instance
func NewMux(lc fx.Lifecycle, logger *zap.Logger, cfg configuration.Configuration) *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix(cfg.HTTP.Endpoint).Subrouter()
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port)

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				logger.Debug("started http server", zap.String("address", addr),
					zap.String("endpoint", cfg.HTTP.Endpoint))
				_ = server.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Debug("stopped http server", zap.String("address", addr),
				zap.String("endpoint", cfg.HTTP.Endpoint))
			return server.Shutdown(ctx)
		},
	})

	return r
}
