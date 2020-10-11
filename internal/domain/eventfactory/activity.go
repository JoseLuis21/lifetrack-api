package eventfactory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Activity is the event factory for aggregate.Activity
type Activity struct{}

// NewActivityCreated returns a pre-build Domain event for aggregate.Activity creation
func (a Activity) NewActivityCreated(ag aggregate.Activity) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "activity",
		Action:        "created",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "activity",
		Body:          adapter.ActivityAdapter{}.ToModel(ag),
		Snapshot:      nil,
	})

	return *e
}

// NewActivityUpdated returns a pre-build Domain event for aggregate.Activity mutations
func (a Activity) NewActivityUpdated(ag, snapshot aggregate.Activity) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "activity",
		Action:        "updated",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "activity",
		Body:          adapter.ActivityAdapter{}.ToModel(ag),
		Snapshot:      adapter.ActivityAdapter{}.ToModel(snapshot),
	})

	return *e
}

// NewActivityRemoved returns a pre-build Domain event for aggregate.Activity removal
func (a Activity) NewActivityRemoved(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
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
func (a Activity) NewActivityRestored(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
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
func (a Activity) NewActivityHardRemoved(ag aggregate.Activity) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "activity",
		Action:        "hard_removed",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "activity",
		Body:          nil,
		Snapshot:      ag,
	})

	return *e
}
