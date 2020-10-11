package eventfactory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Occurrence is the event factory for aggregate.Occurrence
type Occurrence struct{}

// NewOccurrenceCreated returns a pre-build Domain event for Occurrence creation
func (o Occurrence) NewOccurrenceCreated(ag aggregate.Occurrence) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "occurrence",
		Action:        "added",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "Occurrence",
		Body:          adapter.OccurrenceAdapter{}.ToModel(ag),
		Snapshot:      nil,
	})

	return *e
}

// NewOccurrenceUpdated returns a pre-build Domain event for Occurrence mutations
func (o Occurrence) NewOccurrenceUpdated(ag, snapshot aggregate.Occurrence) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "occurrence",
		Action:        "updated",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "Occurrence",
		Body:          adapter.OccurrenceAdapter{}.ToModel(ag),
		Snapshot:      adapter.OccurrenceAdapter{}.ToModel(snapshot),
	})

	return *e
}

// NewOccurrenceRemoved returns a pre-build Domain event for Occurrence removal
func (o Occurrence) NewOccurrenceRemoved(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "occurrence",
		Action:        "removed",
		AggregateID:   id.Get(),
		AggregateName: "Occurrence",
		Body:          nil,
		Snapshot:      nil,
	})

	return *e
}

// NewOccurrenceRestored returns a pre-build Domain event for Occurrence removal
func (o Occurrence) NewOccurrenceRestored(id value.CUID) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "occurrence",
		Action:        "restored",
		AggregateID:   id.Get(),
		AggregateName: "Occurrence",
		Body:          nil,
		Snapshot:      nil,
	})

	return *e
}

// NewOccurrenceHardRemoved returns a pre-build Domain event for Occurrence permanently removal
func (o Occurrence) NewOccurrenceHardRemoved(ag aggregate.Occurrence) event.Domain {
	e, _ := event.NewDomain(event.DomainArgs{
		Service:       "occurrence",
		Action:        "hard_removed",
		AggregateID:   ag.Get().ID.Get(),
		AggregateName: "Occurrence",
		Body:          nil,
		Snapshot:      ag,
	})

	return *e
}
