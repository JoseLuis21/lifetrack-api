package category

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"

	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Get request a single category
type Get struct {
	repo repository.Category
}

// NewGet get a new Get query implementation
func NewGet(r repository.Category) *Get {
	return &Get{
		repo: r,
	}
}

func (q Get) Query(ctx context.Context, id string) (*model.Category, error) {
	idCUID := new(value.CUID)
	if err := idCUID.Set(id); err != nil {
		return nil, err
	}

	c, err := q.repo.FetchByID(ctx, *idCUID)
	if err != nil {
		return nil, err
	}

	return adapter.CategoryAdapter{}.ToModel(*c), nil
}
