package factory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"time"
)

func NewCategory(title, user, description string) (*aggregate.Category, error) {
	titleP, err := value.NewTitle(title)
	if err != nil {
		return nil, err
	}

	desc, err := value.NewDescription(description)
	if err != nil {
		return nil, err
	}

	return &aggregate.Category{
		ID:          value.NewUUID(),
		Title:       titleP,
		Description: desc,
		User:        user,
		CreateTime:  time.Now().UTC(),
		UpdateTime:  time.Now().UTC(),
	}, nil
}
