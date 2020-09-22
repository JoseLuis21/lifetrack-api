package repository

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/domain"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

type Category interface {
	Save(ctx context.Context, c aggregate.Category) error
	FetchByID(ctx context.Context, id value.UUID) (*model.Category, error)
	Fetch(ctx context.Context, token string, limit int64, filter domain.FilterMap) ([]*model.Category, string, error)
}
