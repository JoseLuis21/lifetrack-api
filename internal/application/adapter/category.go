package adapter

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"time"
)

// CategoryAdapter adapts different types of category structs
type CategoryAdapter struct{}

// ToAggregate parses a category aggregate root to a read-only model
func (a CategoryAdapter) ToModel(ag aggregate.Category) *model.Category {
	return &model.Category{
		ID:          ag.GetRoot().ID.Get(),
		Title:       ag.GetRoot().Title.Get(),
		Description: ag.GetRoot().Description.Get(),
		User:        ag.GetRoot().User,
		Theme:       ag.GetRoot().Theme.Get(),
		CreateTime:  ag.GetRoot().Metadata.GetCreateTime().Unix(),
		UpdateTime:  ag.GetRoot().Metadata.GetUpdateTime().Unix(),
		Active:      ag.GetRoot().Metadata.GetState(),
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

	theme, err := value.NewTheme(m.Theme)
	if err != nil {
		return nil, err
	}

	meta := new(value.Metadata)
	_ = meta.SetCreateTime(time.Unix(m.CreateTime, 0).UTC())
	_ = meta.SetUpdateTime(time.Unix(m.UpdateTime, 0).UTC())
	_ = meta.SetState(m.Active)

	ag := new(aggregate.Category)
	ag.SetRoot(&entity.Category{
		ID:          id,
		Title:       titleP,
		Description: desc,
		User:        m.User,
		Theme:       theme,
		Metadata:    meta,
	})

	return ag, nil
}
