package eventbus

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"sync"
)

// InMemory is the event.Bus implementation using in-memory
type InMemory struct {
	cfg infrastructure.Configuration
	mu  *sync.Mutex
}

// NewInMemory creates a concrete struct of InMemory event bus
func NewInMemory(cfg infrastructure.Configuration) *InMemory {
	return &InMemory{
		cfg: cfg,
		mu:  new(sync.Mutex),
	}
}

func (b *InMemory) Publish(ctx context.Context, e ...event.Domain) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(e) == 0 {
		return exception.NewRequiredField("domain event")
	}

	return nil
}

func (b *InMemory) SubscribeTo(ctx context.Context, t event.Topic) (chan *event.Domain, error) {
	return nil, nil
}
