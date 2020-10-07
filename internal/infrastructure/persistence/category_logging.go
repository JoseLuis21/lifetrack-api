package persistence

import (
	"context"
	"time"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"go.uber.org/zap"
)

// CategoryLog is a CoR implementation for repository.Category logging
type CategoryLog struct {
	Log  *zap.Logger
	Next repository.Category
}

func (r CategoryLog) Save(ctx context.Context, c aggregate.Category) (err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to save category",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("action", "save"),
				zap.Duration("backoff", time.Since(init)),
				zap.String("err", err.Error()),
				zap.String("id", c.GetRoot().ID.Get()),
			)
			return
		}

		r.Log.Info("succeed to save category",
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "save"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", c.GetRoot().ID.Get()),
		)
	}(time.Now())

	err = r.Next.Save(ctx, c)
	return
}

func (r CategoryLog) FetchByID(ctx context.Context, id value.CUID) (category *model.Category, err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to fetch category by id",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("action", "fetch_by_id"),
				zap.Duration("backoff", time.Since(init)),
				zap.String("err", err.Error()),
				zap.String("id", id.Get()),
			)
			return
		}

		r.Log.Info("succeed to fetch category by id",
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "fetch_by_id"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", id.Get()),
		)
	}(time.Now())

	category, err = r.Next.FetchByID(ctx, id)
	return
}

func (r CategoryLog) Fetch(ctx context.Context, token string, limit int64, criteria shared.CategoryCriteria) (
	categories []*model.Category, nextToken string, err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to fetch category",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("action", "fetch"),
				zap.Duration("backoff", time.Since(init)),
				zap.String("err", err.Error()),
				zap.String("next_token", token),
				zap.Int64("page_size", limit),
			)
			return
		}

		r.Log.Info("succeed to fetch category",
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "fetch"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("next_token", token),
			zap.Int64("page_size", limit),
		)
	}(time.Now())

	categories, nextToken, err = r.Next.Fetch(ctx, token, limit, criteria)
	return
}

func (r CategoryLog) Replace(ctx context.Context, c aggregate.Category) (err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to replace category",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("action", "replace"),
				zap.Duration("backoff", time.Since(init)),
				zap.String("err", err.Error()),
				zap.String("id", c.GetRoot().ID.Get()),
			)
			return
		}

		r.Log.Info("succeed to replace category",
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "replace"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", c.GetRoot().ID.Get()),
		)
	}(time.Now())

	err = r.Next.Replace(ctx, c)
	return
}

func (r CategoryLog) HardRemove(ctx context.Context, id value.CUID) (err error) {
	defer func(init time.Time) {
		if err != nil {
			r.Log.Error("failed to hard_remove category",
				zap.String("module", "infrastructure.persistence.category"),
				zap.String("action", "hard_remove"),
				zap.Duration("backoff", time.Since(init)),
				zap.String("err", err.Error()),
				zap.String("id", id.Get()),
			)
			return
		}

		r.Log.Info("succeed to hard_remove category",
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "hard_remove"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", id.Get()),
		)
	}(time.Now())

	err = r.Next.HardRemove(ctx, id)
	return
}
