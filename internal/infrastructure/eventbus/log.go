package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"go.uber.org/zap"
)

type logBus struct {
	Logger *zap.Logger
	Next   event.Bus
}

func (l logBus) Publish(ctx context.Context, events ...event.Domain) (err error) {
	defer func(startTime time.Time) {
		fields := []zap.Field{
			zap.Duration("backoff", time.Since(startTime)),
		}
		if eventsJSON, errJ := json.Marshal(events); errJ == nil {
			fields = append(fields, zap.String("events", string(eventsJSON)))
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			l.Logger.Error("failed to publish events", fields...)
			return
		}

		l.Logger.Info("published events to topic", fields...)
	}(time.Now())

	err = l.Next.Publish(ctx, events...)
	return
}

func (l logBus) SubscribeTo(ctx context.Context, topic string, handler event.Handler) (err error) {
	defer func(startTime time.Time) {
		fields := []zap.Field{
			zap.String("topic", topic),
			zap.Duration("backoff", time.Since(startTime)),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			l.Logger.Error("failed to subscribe", fields...)
			return
		}

		l.Logger.Info("subscribed to topic", fields...)
	}(time.Now())

	err = l.Next.SubscribeTo(ctx, topic, handler)
	return
}
