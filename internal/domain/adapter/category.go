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

// ToModel parses a category aggregate root to a read model for queries
func (a CategoryAdapter) ToModel(category aggregate.Category) *model.Category {
	return &model.Category{
		ID:          category.Get().ID.Get(),
		Title:       category.Get().Title.Get(),
		Description: category.Get().Description.Get(),
		User:        category.GetUser(),
		Color:       category.Get().Color.Get(),
		Image:       category.Get().Image.Get(),
		CreateTime:  category.Get().Metadata.GetCreateTime().Unix(),
		UpdateTime:  category.Get().Metadata.GetUpdateTime().Unix(),
		Active:      category.Get().Metadata.GetState(),
	}
}

// BulkToModel parses a whole slice of category aggregates to read models for queries
func (a CategoryAdapter) BulkToModel(categories ...*aggregate.Category) []*model.Category {
	catsModel := make([]*model.Category, 0)

	for _, cAg := range categories {
		c := a.ToModel(*cAg)
		catsModel = append(catsModel, c)
	}

	return catsModel
}

// ToAggregate parses a category read model to aggregate root
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

	image, err := value.NewImage(m.Image)
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
		Image:       image,
		Metadata:    meta,
	})
	if err := ag.AssignUser(m.User); err != nil {
		return nil, err
	}

	return ag, nil
}
