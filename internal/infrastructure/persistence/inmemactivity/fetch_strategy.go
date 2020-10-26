package inmemactivity

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

// fetchStrategy fetching strategy for activity repositories
type fetchStrategy interface {
	Do(ctx context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error)
}
