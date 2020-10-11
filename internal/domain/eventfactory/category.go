package eventfactory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Category is the event factory for aggregate.Category
type Category struct{}

// NewCategoryCreated returns a pre-build Domain event for category creation
func (c Category) NewCategoryCreated(ag aggregate.Category) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "category",
		Action:        "added",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "category",
		Body:          adapter.CategoryAdapter{}.ToModel(ag),
		Snapshot:      nil,
	})

	return *e
}

// NewCategoryUpdated returns a pre-build Domain event for category mutations
func (c Category) NewCategoryUpdated(ag, snapshot aggregate.Category) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "category",
		Action:        "updated",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "category",
		Body:          adapter.CategoryAdapter{}.ToModel(ag),
		Snapshot:      adapter.CategoryAdapter{}.ToModel(snapshot),
	})

	return *e
}

// NewCategoryRemoved returns a pre-build Domain event for category removal
func (c Category) NewCategoryRemoved(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "category",
		Action:        "removed",
		AggregateID:   id.Get(),
		AggregateName: "category",
		Body:          nil,
		Snapshot:      nil,
	})

	return *e
}

// NewCategoryRestored returns a pre-build Domain event for category removal
func (c Category) NewCategoryRestored(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "category",
		Action:        "restored",
		AggregateID:   id.Get(),
		AggregateName: "category",
		Body:          nil,
		Snapshot:      nil,
	})

	return *e
}

// NewCategoryHardRemoved returns a pre-build Domain event for category permanently removal
func (c Category) NewCategoryHardRemoved(ag aggregate.Category) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "category",
		Action:        "hard_removed",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "category",
		Body:          nil,
		Snapshot:      ag,
	})

	return *e
}
