package event

import "context"

// 	Bus connects a module to the LifeTrack ecosystem through domain and
//	integration events
type Bus interface {
	// Publish sends a group of Domain events to all subscribers
	Publish(ctx context.Context, events ...Domain) error
	// SubscribeTo adds a new subscription and a handler to an specific topic
	SubscribeTo(ctx context.Context, topic string, handler Handler) error
}

// Handler function which handles stream events, acts like a consumer
type Handler func(ctx context.Context, e Domain)
