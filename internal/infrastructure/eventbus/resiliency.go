package eventbus

import (
	"context"
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"go.uber.org/zap"
)

type resiliencyBus struct {
	Logger *zap.Logger
	Next   event.Bus
}

func (r resiliencyBus) Publish(ctx context.Context, events ...event.Domain) error {
	b := breaker.New(3, 1, time.Second*5)
	return b.Run(func() error {
		retry := retrier.New(retrier.ConstantBackoff(3, time.Millisecond*100), retrier.WhitelistClassifier{
			exception.NetworkCall,
		})
		attempts := 0
		return retry.RunCtx(ctx, func(ctxR context.Context) error {
			attempts++
			if err := r.Next.Publish(ctxR, events...); err != nil {
				r.Logger.Info("failed to publish event",
					zap.Int("attempt", attempts),
					zap.Duration("backoff", time.Second),
				)
			}
			return nil
		})
	})
}

func (r resiliencyBus) SubscribeTo(ctx context.Context, topic string, handler event.Handler) error {
	b := breaker.New(3, 1, time.Second*5)
	return b.Run(func() error {
		retry := retrier.New(retrier.ConstantBackoff(3, time.Second*10), retrier.WhitelistClassifier{
			exception.NetworkCall,
		})
		attempts := 0
		return retry.RunCtx(ctx, func(ctxR context.Context) error {
			attempts++
			if err := r.Next.SubscribeTo(ctxR, topic, handler); err != nil {
				r.Logger.Info("failed to subscribe to event",
					zap.Int("attempt", attempts),
					zap.String("topic", topic),
					zap.Duration("backoff", time.Second),
				)
			}
			return nil
		})
	})
}
