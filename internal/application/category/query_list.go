package category

import (
	"context"
	"strconv"

	"github.com/alexandria-oss/common-go/exception"
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
	// TODO: Turn token and limit to value objects
	var limitInt int64
	limitInt = 100

	if limit != "" {
		l, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return nil, "", exception.NewFieldFormat("limit", "integer")
		}
		limitInt = l
	}

	return q.repo.Fetch(ctx, token, limitInt, shared.CategoryCriteria{
		User:    filter["user"],
		Query:   filter["query"],
		OrderBy: filter["order"],
	})
}
