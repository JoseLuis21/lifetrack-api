package eventbus

import (
	"context"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
)

type versioningBus struct {
	Cfg  configuration.Configuration
	Next event.Bus
}

func (v versioningBus) Publish(ctx context.Context, events ...event.Domain) error {
	for i, ev := range events {
		ev.Version = v.Cfg.Version
		ev.Stage = v.Cfg.Stage
		events[i] = ev
	}

	return v.Next.Publish(ctx, events...)
}

func (v versioningBus) SubscribeTo(ctx context.Context, topic string, handler event.Handler) error {
	return v.Next.SubscribeTo(ctx, topic, handler)
}
