package event_factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// NewCategoryCreated returns a pre-build Domain event for category creation
func NewCategoryCreated(ag aggregate.Category) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "category",
		Action:        "created",
		AggregateID:   ag.GetRoot().ID.Get(),
		AggregateName: "category",
		Body:          ag,
		Snapshot:      nil,
	})

	return *e
}

// NewCategoryUpdated returns a pre-build Domain event for category mutations
func NewCategoryUpdated(ag aggregate.Category) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "category",
		Action:        "updated",
		AggregateID:   ag.GetRoot().ID.Get(),
		AggregateName: "category",
		Body:          ag,
		Snapshot:      ag,
	})

	return *e
}

// NewCategoryRemoved returns a pre-build Domain event for category removal
func NewCategoryRemoved(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
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
func NewCategoryRestored(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
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
func NewCategoryHardRemoved(ag aggregate.Category) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "category",
		Action:        "hard_removed",
		AggregateID:   ag.GetRoot().ID.Get(),
		AggregateName: "category",
		Body:          nil,
		Snapshot:      ag,
	})

	return *e
}
