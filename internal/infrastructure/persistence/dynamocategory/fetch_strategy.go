package dynamocategory

import (
	"context"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/aggregate"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
)

// fetchStrategy fetching strategy for category repositories
type fetchStrategy interface {
	Do(ctx context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error)
}
