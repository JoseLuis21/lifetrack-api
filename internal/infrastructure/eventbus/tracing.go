package eventbus

import (
	"context"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"go.opencensus.io/trace"
)

type tracingBus struct {
	Next event.Bus
}

func (t tracingBus) Publish(ctx context.Context, events ...event.Domain) (err error) {
	for _, ev := range events {
		ctxS, span := trace.StartSpan(ctx, "publish")
		span.AddAttributes(t.setAttributeFromEvent(ev)...)
		ev.SpanContext = span.SpanContext()
		if err = t.Next.Publish(ctxS, ev); err != nil {
			span.SetStatus(trace.Status{
				Code:    trace.StatusCodeInternal,
				Message: exception.GetDescription(err),
			})
			span.End()
			return err
		}

		span.End()
	}

	return nil
}

func (t tracingBus) SubscribeTo(ctx context.Context, topic string, handler event.Handler) (err error) {
	ctxS, span := trace.StartSpan(ctx, "subscribe")
	span.AddAttributes(trace.StringAttribute("caller", "lifetrack.event_bus"),
		trace.StringAttribute("action", "subscribe"), trace.StringAttribute("topic", topic))
	defer func() {
		if err != nil {
			span.SetStatus(trace.Status{
				Code:    trace.StatusCodeInternal,
				Message: exception.GetDescription(err),
			})
		}
		span.End()
	}()

	return t.Next.SubscribeTo(ctxS, topic, handler)
}

func (t tracingBus) setAttributeFromEvent(ev event.Domain) []trace.Attribute {
	return []trace.Attribute{trace.StringAttribute("caller", "lifetrack.event_bus"),
		trace.StringAttribute("action", ev.Action), trace.StringAttribute("correlation_id", ev.CorrelationID),
		trace.StringAttribute("aggregate_id", ev.AggregateID), trace.StringAttribute("topic", ev.Topic),
		trace.StringAttribute("version", ev.Version), trace.StringAttribute("stage", ev.Stage),
	}
}
