package eventfactory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
)

type Occurrence struct{}

var occurrenceName = "occurrence"

// ActivityOccurred triggered when a new Occurrence has been added to an Activity
func (o Occurrence) ActivityOccurred(occurrence model.Occurrence) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: occurrenceName,
		Action:        create,
		AggregateID:   occurrence.ID,
		Body:          occurrence,
	})
}

// Updated triggered when an Occurrence has been updated
func (o Occurrence) Updated(occurrence model.Occurrence) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: occurrenceName,
		Action:        update,
		AggregateID:   occurrence.ID,
		Body:          occurrence,
	})
}

// HardRemoved triggered when an Occurrence has been removed permanently
func (o Occurrence) HardRemoved(occurrenceID string) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: occurrenceName,
		Action:        hardRemove,
		AggregateID:   occurrenceID,
		Body:          nil,
	})
}
