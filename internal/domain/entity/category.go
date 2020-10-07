package entity

import (
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Category entity used to group n-activities
type Category struct {
	ID          *value.CUID
	Title       *value.Title
	Description *value.Description
	User        string
	Theme       *value.Theme
	Metadata    *value.Metadata
}

// NewCategory creates a category entity receiving primitive-only data
func NewCategory(title, user, description, theme string) (*Category, error) {
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

	return &Category{
		ID:          value.NewCUID(),
		Title:       titleP,
		Description: desc,
		User:        user,
		Theme:       t,
		Metadata:    value.NewMetadata(),
	}, nil
}

// IsValid validate entity data
func (c Category) IsValid() error {
	// - Required: title, id, user

	suffix := "category_"
	if c.ID == nil || c.ID.Get() == "" {
		return exception.NewRequiredField(suffix + "id")
	} else if c.Title == nil || c.Title.Get() == "" {
		return exception.NewRequiredField(suffix + "title")
	} else if c.User == "" {
		return exception.NewRequiredField(suffix + "user")
	}

	return nil
}

// Update mutates editable data and sets UpdateTime metadata to current time in UTC
func (c *Category) Update(title, description, theme string) error {
	if err := c.Title.Set(title); title != "" && err != nil {
		return err
	}
	if err := c.Description.Set(description); description != "" && err != nil {
		return err
	}
	if err := c.Theme.Set(theme); theme != "" && err != nil {
		return err
	}
	_ = c.Metadata.SetUpdateTime(time.Now().UTC())

	if err := c.IsValid(); err != nil {
		return err
	}

	return nil
}

// Remove sets active flag to false and sets UpdateTime metadata to current time in UTC
func (c *Category) Remove() {
	_ = c.Metadata.SetUpdateTime(time.Now().UTC())
	_ = c.Metadata.SetState(false)
}

// Restore set active flag to true and sets UpdateTime metadata to current time in UTC
func (c *Category) Restore() {
	_ = c.Metadata.SetUpdateTime(time.Now().UTC())
	_ = c.Metadata.SetState(true)
}
