package category

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"go.uber.org/zap"
)

// NewCategory wraps an existing category repository with required observability and resiliency
func NewCategory(r repository.Category, logger *zap.Logger) repository.Category {
	// TODO: Add monitoring, distributed tracing and circuit-breaker w/ retry policy patterns
	repo := Log{
		Log:  logger,
		Next: r,
	}

	return repo
}
