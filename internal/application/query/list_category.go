package query

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"strconv"
)

type ListCategories struct {
	repo repository.Category
}

// NewListCategories get list of categories
func NewListCategories(r repository.Category) *ListCategories {
	return &ListCategories{
		repo: r,
	}
}

func (q ListCategories) Query(ctx context.Context, token, limit string) ([]*model.Category, string, error) {
	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, "", exception.NewFieldFormat("limit", "integer")
	}

	return q.repo.Fetch(ctx, token, int(l))
}
