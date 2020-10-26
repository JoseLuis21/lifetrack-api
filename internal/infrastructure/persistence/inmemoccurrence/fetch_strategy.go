package inmemoccurrence

import (
	"context"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/aggregate"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
)

// fetchStrategy fetching strategy for occurrence repositories
type fetchStrategy interface {
	Do(ctx context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error)
}
