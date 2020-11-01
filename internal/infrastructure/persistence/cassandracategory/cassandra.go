package cassandracategory

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/aggregate"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
	"sync"
)

type Repository struct {
	s  *gocql.Session
	mu *sync.RWMutex
}

// Singleton is required in this layer to keep mutex consistency between instances
var (
	singleton     *Repository
	singletonLock = new(sync.Once)
)

func NewRepository(s *gocql.Session) *Repository {
	singletonLock.Do(func() {
		singleton = &Repository{
			s:  s,
			mu: new(sync.RWMutex),
		}
	})

	return singleton
}

func (r Repository) Save(ctx context.Context, category aggregate.Category) error {
	panic("not implemented")
}

func (r Repository) Fetch(ctx context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	panic("not implemented")
}

func (r Repository) Remove(ctx context.Context, id string) error {
	panic("not implemented")
}
