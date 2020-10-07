package repository

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Activity offers an abstract data-link to Activity persistence layer(s)
type Activity interface {
	// Save stores an aggregate.Activity
	Save(ctx context.Context, c aggregate.Activity) error
	// FetchByID returns an aggregate.Activity by its value.UUID
	FetchByID(ctx context.Context, id value.CUID) (*model.Activity, error)
	// Fetch returns an aggregate.Activity slice, accepts multiple filters and params such as nextToken and limit for
	// pagination
	Fetch(ctx context.Context, token string, limit int64, criteria shared.ActivityCriteria) ([]*model.Activity,
		string, error)
	// Replace mutates completely an aggregate.Activity
	Replace(ctx context.Context, c aggregate.Activity) error
	// HardRemove permanently removes an aggregate.Activity
	HardRemove(ctx context.Context, id value.CUID) error
}
