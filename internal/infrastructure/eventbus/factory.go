package eventbus

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"go.uber.org/zap"
)

// NewEventBus wraps an existing event bus with required observability and resiliency
func NewEventBus(b event.Bus, logger *zap.Logger) event.Bus {
	// TODO: Add monitoring, distributed tracing and circuit-breaker w/ retry policy patterns
	return Log{
		Logger: logger,
		Next:   b,
	}
}
