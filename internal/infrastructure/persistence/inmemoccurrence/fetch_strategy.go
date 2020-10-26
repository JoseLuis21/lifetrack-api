package inmemoccurrence

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

// fetchStrategy fetching strategy for occurrence repositories
type fetchStrategy interface {
	Do(ctx context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error)
}
