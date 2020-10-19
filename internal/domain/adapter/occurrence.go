package adapter

import (
	"fmt"
	"time"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// OccurrenceAdapter adapts different types of Occurrence structs
type OccurrenceAdapter struct{}

// ToModel parses a Occurrence aggregate root to a read-only model
func (a OccurrenceAdapter) ToModel(ag aggregate.Occurrence) *model.Occurrence {
	return &model.Occurrence{
		ID:            ag.Get().ID.Get(),
		StartTime:     ag.Get().StartTime.Unix(),
		EndTime:       ag.Get().EndTime.Unix(),
		TotalDuration: ag.Get().TotalDuration.Minutes(),
		Activity:      ag.GetActivity(),
		CreateTime:    ag.Get().Metadata.GetCreateTime().Unix(),
		UpdateTime:    ag.Get().Metadata.GetUpdateTime().Unix(),
		Active:        ag.Get().Metadata.GetState(),
	}
}

// ToAggregate parses a Occurrence read-only model to aggregate root
func (a OccurrenceAdapter) ToAggregate(m model.Occurrence) (*aggregate.Occurrence, error) {
	id := &value.CUID{}
	err := id.Set(m.ID)
	if err != nil {
		return nil, err
	}

	meta := new(value.Metadata)
	meta.SetCreateTime(time.Unix(m.CreateTime, 0).UTC())
	meta.SetUpdateTime(time.Unix(m.UpdateTime, 0).UTC())
	meta.SetState(m.Active)

	// Parse to minutes
	duration, err := time.ParseDuration(fmt.Sprintf("%fm", m.TotalDuration))
	if err != nil {
		return nil, exception.NewFieldFormat("total_duration", "real (R) number")
	}

	ag := new(aggregate.Occurrence)
	ag.Set(&entity.Occurrence{
		ID:            id,
		StartTime:     time.Unix(m.StartTime, 0).UTC(),
		EndTime:       time.Unix(m.StartTime, 0).UTC(),
		TotalDuration: duration,
		Metadata:      meta,
	})
	if err := ag.AssignActivity(m.Activity); err != nil {
		return nil, err
	}

	return ag, nil
}
