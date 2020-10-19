package category

import (
	"context"
	"sync"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// InMemoryRepository in-memory category repository layer
type InMemoryRepository struct {
	items map[string]*aggregate.Category
	mu    *sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: map[string]*aggregate.Category{},
		mu:    new(sync.RWMutex),
	}
}

func (r *InMemoryRepository) Save(_ context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// O(1) performance
	if _, ok := r.items[c.Get().ID.Get()]; ok {
		return exception.NewAlreadyExists("category")
	}

	r.items[c.Get().ID.Get()] = &c
	return nil
}

func (r InMemoryRepository) FetchByID(_ context.Context, id value.CUID) (*aggregate.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.items[id.Get()]
	if !ok {
		return nil, exception.NewNotFound("category")
	}

	return c, nil
}

func (r InMemoryRepository) Fetch(_ context.Context, _ string, limit int64, _ shared.CategoryCriteria) ([]*aggregate.Category, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*aggregate.Category, 0)
	reads := int64(0)
	for _, v := range r.items {
		if reads > limit {
			break
		}

		items = append(items, v)
	}
	if len(items) == 0 {
		return nil, "", exception.NewNotFound("category")
	}

	return items, "", nil
}

func (r *InMemoryRepository) Replace(_ context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[c.Get().ID.Get()]
	if !ok {
		return exception.NewNotFound("category")
	}

	r.items[c.Get().ID.Get()] = &c
	return nil
}

func (r *InMemoryRepository) HardRemove(_ context.Context, id value.CUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[id.Get()]
	if !ok {
		return exception.NewNotFound("category")
	}

	delete(r.items, id.Get())

	return nil
}
