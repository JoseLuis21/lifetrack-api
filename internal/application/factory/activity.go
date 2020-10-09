package factory

import (
	"github.com/neutrinocorp/life-track-api/internal/application/eventfactory"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
)

// NewActivity creates an Activity receiving primitive-only data.
// In addition, this function adds required events to the returning aggregate
func NewActivity(title, category string) (*aggregate.Activity, error) {
	c, err := entity.NewActivity(title, category)
	if err != nil {
		return nil, err
	}

	ag := new(aggregate.Activity)
	ag.SetRoot(c)
	ag.RecordEvent(eventfactory.NewActivityCreated(*ag))

	return ag, nil
}
