package query

import (
	"context"
	"strconv"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
)

// ListCategories request a list of categories
type ListCategories struct {
	repo repository.Category
}

// NewListCategories get list of categories
func NewListCategories(r repository.Category) *ListCategories {
	return &ListCategories{
		repo: r,
	}
}

func (q ListCategories) Query(ctx context.Context, token, limit string, filter map[string]string) ([]*model.Category, string, error) {
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
