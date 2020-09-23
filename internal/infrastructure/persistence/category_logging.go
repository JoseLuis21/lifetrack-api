package persistence

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/domain"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"go.uber.org/zap"
	"time"
)

type CategoryLog struct {
	Log  *zap.Logger
	Next repository.Category
}

func (r CategoryLog) Save(ctx context.Context, c aggregate.Category) (err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to save category",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("err", err.Error()),
				zap.Duration("backoff", time.Since(init)),
				zap.String("title", c.GetRoot().Title.Get()),
				zap.String("user", c.GetRoot().User),
				zap.String("description", c.GetRoot().Description.Get()),
			)
			return
		}

		r.Log.Info("succeed to save category",
			zap.String("module", "infrastructure.persistence.category"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", c.GetRoot().ID.Get()),
			zap.String("title", c.GetRoot().Title.Get()),
			zap.String("user", c.GetRoot().User),
			zap.String("description", c.GetRoot().Description.Get()),
			zap.Time("create_time", c.GetRoot().CreateTime),
		)
	}(time.Now())

	err = r.Next.Save(ctx, c)
	return
}

func (r CategoryLog) FetchByID(ctx context.Context, id value.UUID) (category *model.Category, err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to fetch category by id",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("err", err.Error()),
				zap.Duration("backoff", time.Since(init)),
				zap.String("id", id.Get()),
			)
			return
		}

		r.Log.Info("succeed to fetch category by id",
			zap.String("module", "infrastructure.persistence.category"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", id.Get()),
		)
	}(time.Now())

	category, err = r.Next.FetchByID(ctx, id)
	return
}

func (r CategoryLog) Fetch(ctx context.Context, token string, limit int64, filter domain.FilterMap) (
	categories []*model.Category, nextToken string, err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to fetch category",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("err", err.Error()),
				zap.Duration("backoff", time.Since(init)),
				zap.String("next_token", token),
				zap.Int64("page_size", limit),
			)
			return
		}

		r.Log.Info("succeed to fetch category",
			zap.String("module", "infrastructure.persistence.category"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("next_token", token),
			zap.Int64("page_size", limit),
		)
	}(time.Now())

	categories, nextToken, err = r.Next.Fetch(ctx, token, limit, filter)
	return
}
