package event

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// NewCategoryCreated returns a pre-build Domain event for category creation
func NewCategoryCreated(ag aggregate.Category) Domain {
	e, _ := NewDomain(DomainArgsDTO{
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
func NewCategoryUpdated(ag aggregate.Category) Domain {
	e, _ := NewDomain(DomainArgsDTO{
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
func NewCategoryRemoved(id value.UUID) Domain {
	e, _ := NewDomain(DomainArgsDTO{
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
func NewCategoryRestored(id value.UUID) Domain {
	e, _ := NewDomain(DomainArgsDTO{
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
func NewCategoryHardRemoved(ag aggregate.Category) Domain {
	e, _ := NewDomain(DomainArgsDTO{
		Service:       "category",
		Action:        "hard_removed",
		AggregateID:   ag.GetRoot().ID.Get(),
		AggregateName: "category",
		Body:          nil,
		Snapshot:      ag,
	})

	return *e
}
