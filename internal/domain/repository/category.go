package repository

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
)

// CategoryCriteria sets the Category fetching strategy
type CategoryCriteria struct {
	ID   string
	User string
	Name string
	// Keyword fetch Category containing the given keyword in name and description fields
	Keyword string
	Limit   int64
	Token   string
}

// Category handles aggregate.Category persistence
type Category interface {
	Save(ctx context.Context, category aggregate.Category) error
	Fetch(ctx context.Context, criteria CategoryCriteria) ([]*aggregate.Category, string, error)
	Remove(ctx context.Context, id string) error
}
