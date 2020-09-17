package query

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

type GetCategory struct {
	repo repository.Category
}

// NewGetCategory get a new get category query
func NewGetCategory(r repository.Category) *GetCategory {
	return &GetCategory{
		repo: r,
	}
}

func (q GetCategory) Query(ctx context.Context, id string) (*model.Category, error) {
	idUUID := new(value.UUID)
	if err := idUUID.Set(id); err != nil {
		return nil, err
	}

	return q.repo.FetchByID(ctx, idUUID)
}
