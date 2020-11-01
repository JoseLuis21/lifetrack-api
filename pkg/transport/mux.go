package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/middleware"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/observability"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewMux creates a preconfigured mux.Router, eventually routes must be injected to the instance
func NewMux(lc fx.Lifecycle, logger *zap.Logger, cfg configuration.Configuration) *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix(cfg.HTTP.Endpoint).Subrouter()
	middleware.InjectHTTP(r)
	addr := getHTTPAddress(cfg)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      &ochttp.Handler{Handler: r},
	}
	observability.NewHTTP(logger, cfg)

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

	return r
}
