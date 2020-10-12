package category

import (
	"context"
	"time"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"go.uber.org/zap"
)

// Log is a CoR implementation for repository.Category logging
type Log struct {
	Log  *zap.Logger
	Next repository.Category
}

func (r Log) Save(ctx context.Context, c aggregate.Category) (err error) {
	defer func(init time.Time) {
		fields := []zap.Field{
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "save"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", c.Get().ID.Get()),
			zap.String("title", c.Get().Title.Get()),
			zap.String("user", c.GetUser()),
		}

		if err != nil {
			fields = append(fields, zap.String("err", exception.GetDescription(err)))
			r.Log.Error("failed to save category", fields...)
			return
		}

		r.Log.Info("succeed to save category", fields...)
	}(time.Now())

	err = r.Next.Save(ctx, c)
	return
}

func (r Log) FetchByID(ctx context.Context, id value.CUID) (category *model.Category, err error) {
	defer func(init time.Time) {
		fields := []zap.Field{
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "fetch_by_id"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", id.Get()),
		}

		if err != nil {
			fields = append(fields, zap.String("err", exception.GetDescription(err)))
			r.Log.Error("failed to fetch category by id", fields...)
			return
		}

		r.Log.Info("succeed to fetch category by id", fields...)
	}(time.Now())

	category, err = r.Next.FetchByID(ctx, id)
	return
}

func (r Log) Fetch(ctx context.Context, token string, limit int64, criteria shared.CategoryCriteria) (
	categories []*model.Category, nextToken string, err error) {
	defer func(init time.Time) {
		fields := []zap.Field{
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "fetch"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("next_page", token),
			zap.Int64("page_size", limit),
			zap.String("user", criteria.User),
			zap.String("query", criteria.Query),
			zap.String("order_by", criteria.OrderBy),
		}

		if err != nil {
			fields = append(fields, zap.String("err", exception.GetDescription(err)))
			r.Log.Error("failed to fetch category", fields...)
			return
		}

		r.Log.Info("succeed to fetch category", fields...)
	}(time.Now())

	categories, nextToken, err = r.Next.Fetch(ctx, token, limit, criteria)
	return
}

func (r Log) Replace(ctx context.Context, c aggregate.Category) (err error) {
	defer func(init time.Time) {
		fields := []zap.Field{
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "replace"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", c.Get().ID.Get()),
			zap.String("title", c.Get().Title.Get()),
			zap.String("user", c.GetUser()),
			zap.Bool("state", c.Get().Metadata.GetState()),
		}

		if err != nil {
			fields = append(fields, zap.String("err", exception.GetDescription(err)))
			r.Log.Error("failed to replace category", fields...)
			return
		}

		r.Log.Info("succeed to replace category", fields...)
	}(time.Now())

	err = r.Next.Replace(ctx, c)
	return
}

func (r Log) HardRemove(ctx context.Context, id value.CUID) (err error) {
	defer func(init time.Time) {
		fields := []zap.Field{
			zap.String("module", "infrastructure.persistence.category"),
			zap.String("action", "hard_remove"),
			zap.Duration("backoff", time.Since(init)),
			zap.String("id", id.Get()),
		}

		if err != nil {
			fields = append(fields, zap.String("err", exception.GetDescription(err)))
			r.Log.Error("failed to hard_remove category", fields...)
			return
		}

		r.Log.Info("succeed to hard_remove category", fields...)
	}(time.Now())

	err = r.Next.HardRemove(ctx, id)
	return
}
