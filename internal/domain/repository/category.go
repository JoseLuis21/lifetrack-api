package repository

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Category offers an abstract data-link to Category persistence layer(s)
type Category interface {
	// Save stores an aggregate.Category
	Save(ctx context.Context, c aggregate.Category) error
	// FetchByID returns an aggregate.Category by its value.UUID
	FetchByID(ctx context.Context, id value.CUID) (*model.Category, error)
	// Fetch returns an aggregate.Category slice, accepts multiple filters and params such as nextToken and limit for
	// pagination
	Fetch(ctx context.Context, token string, limit int64, criteria shared.CategoryCriteria) ([]*model.Category, string, error)
	// Replace mutates completely an aggregate.Category
	Replace(ctx context.Context, c aggregate.Category) error
	// HardRemove permanently removes an aggregate.Category
	HardRemove(ctx context.Context, id value.CUID) error
}
