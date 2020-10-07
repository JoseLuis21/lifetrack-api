package adapter

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"time"
)

// ActivityAdapter adapts different types of Activity structs
type ActivityAdapter struct{}

// ToAggregate parses a Activity aggregate root to a read-only model
func (a ActivityAdapter) ToModel(ag aggregate.Activity) *model.Activity {
	return &model.Activity{
		ID:         ag.GetRoot().ID.Get(),
		Title:      ag.GetRoot().Title.Get(),
		Category:   ag.GetRoot().Category.Get(),
		CreateTime: ag.GetRoot().Metadata.GetCreateTime().Unix(),
		UpdateTime: ag.GetRoot().Metadata.GetUpdateTime().Unix(),
		Active:     ag.GetRoot().Metadata.GetState(),
	}
}

// ToAggregate parses a Activity read-only model to aggregate root
func (a ActivityAdapter) ToAggregate(m model.Activity) (*aggregate.Activity, error) {
	id := &value.CUID{}
	err := id.Set(m.ID)
	if err != nil {
		return nil, err
	}

	titleP, err := value.NewTitle("Activity_title", m.Title)
	if err != nil {
		return nil, err
	}

	category := &value.CUID{}
	err = id.Set(m.ID)
	if err != nil {
		return nil, err
	}

	meta := new(value.Metadata)
	_ = meta.SetCreateTime(time.Unix(m.CreateTime, 0).UTC())
	_ = meta.SetUpdateTime(time.Unix(m.UpdateTime, 0).UTC())
	_ = meta.SetState(m.Active)

	ag := new(aggregate.Activity)
	ag.SetRoot(&entity.Activity{
		ID:       id,
		Title:    titleP,
		Category: category,
		Metadata: meta,
	})

	return ag, nil
}
