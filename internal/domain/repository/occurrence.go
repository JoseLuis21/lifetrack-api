package repository

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Occurrence offers an abstract data-link to Occurrence persistence layer(s)
type Occurrence interface {
	// Save stores an aggregate.Occurrence
	Save(ctx context.Context, c aggregate.Occurrence) error
	// FetchByID returns an aggregate.Occurrence by its value.UUID
	FetchByID(ctx context.Context, id value.CUID) (*model.Occurrence, error)
	// Fetch returns an aggregate.Occurrence slice, accepts multiple filters and params such as nextToken and limit for
	// pagination
	Fetch(ctx context.Context, token string, limit int64, criteria shared.OccurrenceCriteria) ([]*model.Occurrence,
		string, error)
	// Replace mutates completely an aggregate.Occurrence
	Replace(ctx context.Context, c aggregate.Occurrence) error
	// HardRemove permanently removes an aggregate.Occurrence
	HardRemove(ctx context.Context, id value.CUID) error
}
