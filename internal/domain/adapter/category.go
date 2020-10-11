package adapter

import (
	"time"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// CategoryAdapter adapts different types of category structs
type CategoryAdapter struct{}

// ToModel parses a category aggregate root to a read-only model
func (a CategoryAdapter) ToModel(ag aggregate.Category) *model.Category {
	return &model.Category{
		ID:          ag.Get().ID.Get(),
		Title:       ag.Get().Title.Get(),
		Description: ag.Get().Description.Get(),
		User:        ag.GetUser(),
		Color:       ag.Get().Color.Get(),
		CreateTime:  ag.Get().Metadata.GetCreateTime().Unix(),
		UpdateTime:  ag.Get().Metadata.GetUpdateTime().Unix(),
		Active:      ag.Get().Metadata.GetState(),
	}
}

// ToAggregate parses a category read-only model to aggregate root
func (a CategoryAdapter) ToAggregate(m model.Category) (*aggregate.Category, error) {
	id := &value.CUID{}
	err := id.Set(m.ID)
	if err != nil {
		return nil, err
	}

	titleP, err := value.NewTitle("category_title", m.Title)
	if err != nil {
		return nil, err
	}

	desc, err := value.NewDescription("category_description", m.Description)
	if err != nil {
		return nil, err
	}

	color, err := value.NewColor(m.Color)
	if err != nil {
		return nil, err
	}

	meta := new(value.Metadata)
	meta.SetCreateTime(time.Unix(m.CreateTime, 0).UTC())
	meta.SetUpdateTime(time.Unix(m.UpdateTime, 0).UTC())
	meta.SetState(m.Active)

	ag := new(aggregate.Category)
	ag.Set(&entity.Category{
		ID:          id,
		Title:       titleP,
		Description: desc,
		Color:       color,
		Metadata:    meta,
	})
	if err := ag.AssignUser(m.User); err != nil {
		return nil, err
	}

	return ag, nil
}
