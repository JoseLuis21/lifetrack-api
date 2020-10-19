package category

import (
	"context"
	"strconv"

	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"

	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
)

// List request a list of categories
type List struct {
	repo repository.Category
}

// NewList get a new List query implementation
func NewList(r repository.Category) *List {
	return &List{
		repo: r,
	}
}

func (q List) Query(ctx context.Context, token, limit string, filter map[string]string) ([]*model.Category, string, error) {
	c, next, err := q.repo.Fetch(ctx, token, q.sanitizeLimit(limit), q.generateCriteria(filter))
	if err != nil {
		return nil, "", err
	}

	return adapter.CategoryAdapter{}.BulkToModel(c...), next, nil
}

// sanitizeLimit clean corrupted limit values, returns 10 as default value
func (q List) sanitizeLimit(limit string) int64 {
	var limitInt int64
	limitInt = 10

	if limit != "" {
		l, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return limitInt
		}
		limitInt = l
	}

	return limitInt
}

// generateCriteria create a new category criteria struct
func (q List) generateCriteria(filter map[string]string) shared.CategoryCriteria {
	return shared.CategoryCriteria{
		User:    filter["user"],
		Query:   filter["query"],
		OrderBy: filter["order"],
	}
}
