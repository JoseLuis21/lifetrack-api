package repository

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
)

type Category interface {
	Save(ctx context.Context, c *aggregate.Category) error
	FetchByID(ctx context.Context, id string) (*model.Category, error)
}
