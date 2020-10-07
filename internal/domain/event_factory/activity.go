package event_factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// NewActivityCreated returns a pre-build Domain event for aggregate.Activity creation
func NewActivityCreated(ag aggregate.Activity) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "activity",
		Action:        "created",
		AggregateID:   ag.GetRoot().ID.Get(),
		AggregateName: "activity",
		Body:          ag,
		Snapshot:      nil,
	})

	return *e
}

// NewActivityUpdated returns a pre-build Domain event for aggregate.Activity mutations
func NewActivityUpdated(ag aggregate.Activity) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "activity",
		Action:        "updated",
		AggregateID:   ag.GetRoot().ID.Get(),
		AggregateName: "activity",
		Body:          ag,
		Snapshot:      ag,
	})

	return *e
}

// NewActivityRemoved returns a pre-build Domain event for aggregate.Activity removal
func NewActivityRemoved(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "activity",
		Action:        "removed",
		AggregateID:   id.Get(),
		AggregateName: "activity",
		Body:          nil,
		Snapshot:      nil,
	})

	return *e
}

// NewActivityRestored returns a pre-build Domain event for aggregate.Activity removal
func NewActivityRestored(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "activity",
		Action:        "restored",
		AggregateID:   id.Get(),
		AggregateName: "activity",
		Body:          nil,
		Snapshot:      nil,
	})

	return *e
}

// NewActivityHardRemoved returns a pre-build Domain event for aggregate.Activity permanently removal
func NewActivityHardRemoved(ag aggregate.Activity) event.Domain {
	e, _ := event.NewDomain(event.DomainArgsDTO{
		Service:       "activity",
		Action:        "hard_removed",
		AggregateID:   ag.GetRoot().ID.Get(),
		AggregateName: "activity",
		Body:          nil,
		Snapshot:      ag,
	})

	return *e
}
