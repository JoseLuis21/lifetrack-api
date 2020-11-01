package eventbus

import (
	"context"
	"sync"

	"github.com/asaskevich/EventBus"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
)

var (
	inMemorySingleton *InMemory
	inMemoryLock      = new(sync.Once)
)

// InMemory event.Bus in memory concrete implementation
type InMemory struct {
	bus EventBus.Bus
	cfg configuration.Configuration
	mu  *sync.Mutex
}

// NewInMemory creates a new InMemory bus
func NewInMemory(c configuration.Configuration) *InMemory {
	inMemoryLock.Do(func() {
		inMemorySingleton = &InMemory{bus: EventBus.New(), cfg: c, mu: new(sync.Mutex)}
	})
	return inMemorySingleton
}

// Publish sends a group of Domain events to all subscribers
func (b *InMemory) Publish(_ context.Context, events ...event.Domain) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, e := range events {
		eJSON, err := e.MarshalBinary()
		if err != nil {
			return err
		}
		b.bus.Publish(e.Topic, eJSON)
	}
	return nil
}

// SubscribeTo adds a new subscription to an specific topic, returns a channel of Domain if exists
func (b *InMemory) SubscribeTo(ctx context.Context, topic string, handler event.Handler) error {
	err := make(chan error)
	ctxC, cancel := context.WithCancel(ctx)
	go func() {
		if errS := b.bus.SubscribeAsync(topic, func(incomeEvent []byte) {
			e := new(event.Domain)
			if errM := e.UnmarshalBinary(incomeEvent); errM != nil {
				err <- errM
				return
			}

			handler(ctxC, *e)
		}, false); err != nil {
			err <- errS
			cancel()
		}
	}()

	return <-err
}
