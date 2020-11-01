package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"

	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/resiliency"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/tracing"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewMux creates a preconfigured mux.Router, eventually routes must be injected to the instance
func NewMux(lc fx.Lifecycle, logger *zap.Logger, cfg configuration.Configuration) *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix(cfg.HTTP.Endpoint).Subrouter()
	setMiddlewares(r)
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      &ochttp.Handler{Handler: r},
	}
	setObservability(logger, cfg)

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

// setObservability handles observability depending on the current development stage, strategy pattern applied
func setObservability(logger *zap.Logger, cfg configuration.Configuration) {
	//	rules
	//	a.	if stage == dev or test, then use jaeger and prometheus OpenCensus exporters
	//	b.	if stage != dev or test, then use production config (AWS X-Ray and CloudWatch)
	switch {
	case cfg.IsDevEnv() || cfg.IsTestEnv():
		tracing.NewJaegerHTTP(logger, cfg)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		// TODO: Add Prometheus
	default:
		tracing.NewJaegerHTTP(logger, cfg)
		// TODO: Add AWS X-Ray
		// 1: 1,000 samples exported
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1 / 1000.0)})
		break
	}
}

func setMiddlewares(router *mux.Router) {
	// router.Use(muxhandlers.RecoveryHandler(muxhandlers.RecoveryLogger(logger)))
	router.Use(muxhandlers.CORS(
		muxhandlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		}),
		muxhandlers.AllowedOrigins([]string{"*"}),
	))
	router.Use(muxhandlers.CompressHandler)
	router.Use(resiliency.HTTPRateLimit)
}
