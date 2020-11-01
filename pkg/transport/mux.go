package transport

import (
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/middleware"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/observability"
	"go.uber.org/zap"
)

// NewMux creates a preconfigured mux.Router, eventually routes must be injected to the instance
func NewMux(logger *zap.Logger, cfg configuration.Configuration) *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix(cfg.HTTP.Endpoint).Subrouter()
	middleware.InjectHTTP(r)
	observability.NewHTTP(logger, cfg)
	return r
}
