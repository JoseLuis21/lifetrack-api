package factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/eventfactory"
)

// NewCategory creates a category receiving primitive-only data.
// In addition, this function adds required events to the returning aggregate
func NewCategory(title, user, description, theme string) (*aggregate.Category, error) {
	c, err := entity.NewCategory(title, description, theme)
	if err != nil {
		return nil, err
	}

	ag := new(aggregate.Category)
	ag.Set(c)
	if err := ag.AssignUser(user); err != nil {
		return nil, err
	}

	if err := ag.IsValid(); err != nil {
		return nil, err
	}

	ag.RecordEvent(eventfactory.Category{}.NewCategoryCreated(*ag))

	return ag, nil
}
