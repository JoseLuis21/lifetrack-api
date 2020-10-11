package adapter

import (
	"time"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// ActivityAdapter adapts different types of Activity structs
type ActivityAdapter struct{}

// ToModel parses a Activity aggregate root to a read-only model
func (a ActivityAdapter) ToModel(ag aggregate.Activity) *model.Activity {
	return &model.Activity{
		ID:            ag.Get().ID.Get(),
		Title:         ag.Get().Title.Get(),
		Category:      ag.GetCategory(),
		AppointedTime: int64(ag.Get().AppointedTime.Get().Minutes()),
		CreateTime:    ag.Get().Metadata.GetCreateTime().Unix(),
		UpdateTime:    ag.Get().Metadata.GetUpdateTime().Unix(),
		Active:        ag.Get().Metadata.GetState(),
	}
}

// ToAggregate parses a Activity read-only model to aggregate root
func (a ActivityAdapter) ToAggregate(m model.Activity) (*aggregate.Activity, error) {
	id := &value.CUID{}
	err := id.Set(m.ID)
	if err != nil {
		return nil, err
	}

	titleP, err := value.NewTitle("activity_title", m.Title)
	if err != nil {
		return nil, err
	}

	meta := new(value.Metadata)
	meta.SetCreateTime(time.Unix(m.CreateTime, 0).UTC())
	meta.SetUpdateTime(time.Unix(m.UpdateTime, 0).UTC())
	meta.SetState(m.Active)

	apTime := new(value.AppointedTime)
	if err := apTime.Set(int(m.AppointedTime)); err != nil {
		return nil, err
	}

	ag := new(aggregate.Activity)
	ag.Set(&entity.Activity{
		ID:            id,
		Title:         titleP,
		AppointedTime: apTime,
		Metadata:      meta,
	})
	if err := ag.AssignCategory(m.Category); err != nil {
		return nil, err
	}

	return ag, nil
}
