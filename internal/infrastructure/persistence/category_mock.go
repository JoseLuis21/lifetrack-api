package persistence

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/application/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"sync"
)

// CategoryMock mocks category repository layer
type CategoryMock struct {
	items map[string]*model.Category
	mu    *sync.RWMutex
}

func NewCategoryMock() *CategoryMock {
	return &CategoryMock{
		items: map[string]*model.Category{},
		mu:    new(sync.RWMutex),
	}
}

func (r *CategoryMock) Save(_ context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// O(1) performance
	if _, ok := r.items[c.GetRoot().ID.Get()]; ok {
		return exception.NewAlreadyExists("category")
	}

	r.items[c.GetRoot().ID.Get()] = adapter.CategoryAdapter{}.ToModel(c)
	return nil
}

func (r CategoryMock) FetchByID(_ context.Context, id value.CUID) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.items[id.Get()]
	if !ok {
		return nil, exception.NewNotFound("category")
	}

	return c, nil
}

func (r CategoryMock) Fetch(_ context.Context, _ string, limit int64, _ shared.CategoryCriteria) ([]*model.Category, string, error) {
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

	return items, "", nil
}

func (r *CategoryMock) Replace(_ context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[c.GetRoot().ID.Get()]
	if !ok {
		return exception.NewNotFound("category")
	}

	r.items[c.GetRoot().ID.Get()] = adapter.CategoryAdapter{}.ToModel(c)
	return nil
}

func (r *CategoryMock) HardRemove(_ context.Context, id value.CUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.items[id.Get()]
	if !ok {
		return exception.NewNotFound("category")
	}

	delete(r.items, id.Get())

	return nil
}
