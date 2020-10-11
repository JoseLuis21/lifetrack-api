package factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/eventfactory"
)

// NewActivity creates an Activity receiving primitive-only data.
// In addition, this function adds required events to the returning aggregate
func NewActivity(title string, appointedTime int) (*aggregate.Activity, error) {
	a, err := entity.NewActivity(title, appointedTime)
	if err != nil {
		return nil, err
	}

	ag := new(aggregate.Activity)
	ag.Set(a)
	ag.RecordEvent(eventfactory.Activity{}.NewActivityCreated(*ag))

	return ag, nil
}
