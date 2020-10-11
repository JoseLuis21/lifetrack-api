package factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/eventfactory"
)

// NewOccurrence creates an Occurrence receiving primitive-only data.
// In addition, this function adds required events to the returning aggregate
func NewOccurrence(startTime, endTime string) (*aggregate.Occurrence, error) {
	o, err := entity.NewOccurrence(startTime, endTime)
	if err != nil {
		return nil, err
	}

	ag := new(aggregate.Occurrence)
	ag.Set(o)
	ag.RecordEvent(eventfactory.Occurrence{}.NewOccurrenceCreated(*ag))

	return ag, nil
}
