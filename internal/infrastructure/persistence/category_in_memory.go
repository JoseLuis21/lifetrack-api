package persistence

import (
	"context"
	"sync"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// CategoryInMemory in-memory category repository layer
type CategoryInMemory struct {
	items map[string]*model.Category
	mu    *sync.RWMutex
}

func NewCategoryInMemory() *CategoryInMemory {
	return &CategoryInMemory{
		items: map[string]*model.Category{},
		mu:    new(sync.RWMutex),
	}
}

func (r *CategoryInMemory) Save(_ context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// O(1) performance
	if _, ok := r.items[c.GetRoot().ID.Get()]; ok {
		return exception.NewAlreadyExists("category")
	}

	r.items[c.GetRoot().ID.Get()] = adapter.CategoryAdapter{}.ToModel(c)
	return nil
}

func (r CategoryInMemory) FetchByID(_ context.Context, id value.CUID) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.items[id.Get()]
	if !ok {
		return nil, exception.NewNotFound("category")
	}

	return c, nil
}

func (r CategoryInMemory) Fetch(_ context.Context, _ string, limit int64, _ shared.CategoryCriteria) ([]*model.Category, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*model.Category, 0)
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

func (r *CategoryInMemory) Replace(_ context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[c.GetRoot().ID.Get()]
	if !ok {
		return exception.NewNotFound("category")
	}

	r.items[c.GetRoot().ID.Get()] = adapter.CategoryAdapter{}.ToModel(c)
	return nil
}

func (r *CategoryInMemory) HardRemove(_ context.Context, id value.CUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[id.Get()]
	if !ok {
		return exception.NewNotFound("category")
	}

	delete(r.items, id.Get())

	return nil
}
