package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/aggregate"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
	"go.uber.org/zap"
)

type categoryLog struct {
	Logger *zap.Logger
	Next   repository.Category
}

func (c categoryLog) Save(ctx context.Context, category aggregate.Category) (err error) {
	defer func(startTime time.Time) {
		fields := []zap.Field{
			zap.String("aggregate_id", category.ID()),
			zap.Duration("backoff", time.Since(startTime)),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			c.Logger.Error("failed to create category", fields...)
			return
		}

		c.Logger.Info("category created", fields...)
	}(time.Now())

	err = c.Next.Save(ctx, category)
	return
}

func (c categoryLog) Fetch(ctx context.Context, criteria repository.CategoryCriteria) (categories []*aggregate.Category,
	nextToken string, err error) {
	defer func(startTime time.Time) {
		fields := []zap.Field{
			zap.String("aggregate_id", criteria.ID),
			zap.String("keyword", criteria.Keyword),
			zap.String("name", criteria.Name),
			zap.String("token", criteria.Token),
			zap.String("user", criteria.User),
			zap.Int64("limit", criteria.Limit),
			zap.Duration("backoff", time.Since(startTime)),
		}

		if err != nil && !errors.Is(err, exception.NotFound) {
			fields = append(fields, zap.Error(err))
			c.Logger.Error("failed to fetch category", fields...)
			return
		} else if err != nil {
			fields = append(fields, zap.Error(err))
			c.Logger.Warn("category not found", fields...)
			return
		}

		c.Logger.Info("category fetch", fields...)
	}(time.Now())

	categories, nextToken, err = c.Next.Fetch(ctx, criteria)
	return
}

func (c categoryLog) Remove(ctx context.Context, id string) (err error) {
	defer func(startTime time.Time) {
		fields := []zap.Field{
			zap.String("aggregate_id", id),
			zap.Duration("backoff", time.Since(startTime)),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			c.Logger.Error("failed to remove category", fields...)
			return
		}

		c.Logger.Info("category removed", fields...)
	}(time.Now())

	err = c.Next.Remove(ctx, id)
	return
}
