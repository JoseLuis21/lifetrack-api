package factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event_factory"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// NewCategory creates a category receiving primitive-only data.
// In addition, this function adds required events to the returning aggregate
func NewCategory(title, user, description, theme string) (*aggregate.Category, error) {
	titleP, err := value.NewTitle("category_title", title)
	if err != nil {
		return nil, err
	}

	desc, err := value.NewDescription("category_description", description)
	if err != nil {
		return nil, err
	}

	t, err := value.NewTheme(theme)
	if err != nil {
		return nil, err
	}

	ag := new(aggregate.Category)
	ag.SetRoot(&entity.Category{
		ID:          value.NewCUID(),
		Title:       titleP,
		Description: desc,
		User:        user,
		Theme:       t,
		Metadata:    value.NewMetadata(),
	})
	ag.RecordEvent(event_factory.NewCategoryCreated(*ag))

	return ag, nil
}
