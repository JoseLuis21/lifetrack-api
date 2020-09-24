package factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"time"
)

// NewCategory creates a category receiving primitive-only data.
// In addition, this function adds required events to the returning aggregate
func NewCategory(title, user, description string) (*aggregate.Category, error) {
	titleP, err := value.NewTitle(title)
	if err != nil {
		return nil, err
	}

	desc, err := value.NewDescription(description)
	if err != nil {
		return nil, err
	}

	ag := new(aggregate.Category)
	ag.SetRoot(&entity.Category{
		ID:          value.NewUUID(),
		Title:       titleP,
		Description: desc,
		User:        user,
		CreateTime:  time.Now().UTC(),
		UpdateTime:  time.Now().UTC(),
		Active:      true,
	})
	ag.RecordEvent(event.NewCategoryCreated(*ag))

	return ag, nil
}

// NewCategoryFromModel creates a category receiving read-only model data
func NewCategoryFromModel(p model.Category) (*aggregate.Category, error) {
	id := &value.UUID{}
	err := id.Set(p.ID)
	if err != nil {
		return nil, err
	}

	titleP, err := value.NewTitle(p.Title)
	if err != nil {
		return nil, err
	}

	desc, err := value.NewDescription(p.Description)
	if err != nil {
		return nil, err
	}

	ag := new(aggregate.Category)
	ag.SetRoot(&entity.Category{
		ID:          id,
		Title:       titleP,
		Description: desc,
		User:        p.User,
		CreateTime:  time.Unix(p.CreateTime, 0).UTC(),
		UpdateTime:  time.Unix(p.UpdateTime, 0).UTC(),
		Active:      p.Active,
	})

	return ag, nil
}
