package inmemactivity

import (
	"context"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/aggregate"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
)

// fetchStrategy fetching strategy for activity repositories
type fetchStrategy interface {
	Do(ctx context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error)
}
