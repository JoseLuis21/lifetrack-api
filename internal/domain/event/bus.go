package event

import "context"

// TODO: Add resiliency (retry, circuit breaker) and observability (monitoring, logging, distributed tracing) using CoR pattern

// Bus Event Bus abstraction, represents a message broker
type Bus interface {
	// Publish produces and push an domain event into the Bus
	Publish(ctx context.Context, e ...Domain) error
	// SubscribeTo consumes asynchronously domain events from the Bus
	SubscribeTo(ctx context.Context, t Topic) error
}
