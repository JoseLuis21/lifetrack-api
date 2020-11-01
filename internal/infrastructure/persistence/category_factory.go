package persistence

import (
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
	"go.uber.org/zap"
)

// NewCategory returns a repository.Category with observability and caching strategies
//	Observability: Monitoring, Logging and Tracing
func NewCategory(r repository.Category, logger *zap.Logger) repository.Category {
	return categoryLog{
		Logger: logger,
		Next:   r,
	}
}
