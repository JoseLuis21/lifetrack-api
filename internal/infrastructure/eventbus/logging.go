package eventbus

import (
	"context"
	"time"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"go.uber.org/zap"
)

// Log is a CoR implementation for event.Bus logging
type Log struct {
	Logger *zap.Logger
	Next   event.Bus
}

func (l Log) Publish(ctx context.Context, e ...event.Domain) (err error) {
	defer func(init time.Time) {
		fields := []zap.Field{
			zap.String("module", "infrastructure.eventbus"),
			zap.String("action", "publish"),
			zap.Duration("backoff", time.Since(init)),
		}
		fields = append(fields, zap.Reflect("event", e))

		if err != nil {
			fields = append(fields, zap.String("err", exception.GetDescription(err)))
			l.Logger.Error("failed to publish event", fields...)
			return
		}

		l.Logger.Info("succeed to publish event", fields...)
	}(time.Now())

	err = l.Next.Publish(ctx, e...)
	return
}

func (l Log) SubscribeTo(ctx context.Context, t event.Topic) (chan *event.Domain, error) {
	l.Logger.Info("subscribing to event", []zap.Field{
		zap.String("module", "infrastructure.eventbus"),
		zap.String("action", "subscribe"),
		zap.String("topic", string(t)),
	}...)

	return l.Next.SubscribeTo(ctx, t)
}
