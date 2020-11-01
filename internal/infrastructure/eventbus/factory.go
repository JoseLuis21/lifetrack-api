package eventbus

import (
	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"go.uber.org/zap"
)

// New returns an event.Bus with versioning, observability and resiliency
//	Observability: Monitoring, Logging and Tracing
//	Resiliency: Retry, Circuit Breaker
func New(b event.Bus, cfg configuration.Configuration, logger *zap.Logger) event.Bus {
	// Chaining hierarchy
	// versioning -> logging -> tracing -> resiliency -> monitoring

	// tracing will separate each event to create individual spans
	return versioningBus{
		Cfg: cfg,
		Next: logBus{
			Logger: logger,
			Next: tracingBus{
				Next: resiliencyBus{
					Logger: logger,
					Next:   b,
				},
			},
		},
	}
}
