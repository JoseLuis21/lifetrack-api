package entity

import (
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Category is a group of n-activities tagged by a title
type Category struct {
	ID          *value.CUID
	Title       *value.Title
	Description *value.Description
	Image       *value.Image
	Color       *value.Color
	Metadata    *value.Metadata
}

// NewCategory creates a category entity receiving primitive-only data
func NewCategory(title, description, color, image string) (*Category, error) {
	titleP, err := value.NewTitle("category_title", title)
	if err != nil {
		return nil, err
	}

	desc, err := value.NewDescription("category_description", description)
	if err != nil {
		return nil, err
	}

	cl, err := value.NewColor(color)
	if err != nil {
		return nil, err
	}

	img, err := value.NewImage(image)
	if err != nil {
		return nil, err
	}

	c := &Category{
		ID:          value.NewCUID(),
		Title:       titleP,
		Description: desc,
		Color:       cl,
		Image:       img,
		Metadata:    value.NewMetadata(),
	}

	if err := c.IsValid(); err != nil {
		return nil, err
	}

	return c, nil
}

// IsValid validate entity data
func (c Category) IsValid() error {
	// - Required: title, id
	suffix := "category_"
	if c.ID == nil || c.ID.Get() == "" {
		return exception.NewRequiredField(suffix + "id")
	} else if c.Title == nil || c.Title.Get() == "" {
		return exception.NewRequiredField(suffix + "title")
	}

	return nil
}

// Update mutates data atomically and sets UpdateTime metadata to current time in UTC
func (c *Category) Update(title, description, color, image string) error {
	if err := c.Title.Set(title); title != "" && err != nil {
		return err
	}
	if err := c.Description.Set(description); description != "" && err != nil {
		return err
	}
	if err := c.Color.Set(color); color != "" && err != nil {
		return err
	}
	if err := c.Image.Set(image); image != "" && err != nil {
		return err
	}
	c.Metadata.SetUpdateTime(time.Now().UTC())

	if err := c.IsValid(); err != nil {
		return err
	}

	return nil
}

// Remove sets active flag to false and sets UpdateTime metadata to current time in UTC
func (c *Category) Remove() {
	c.Metadata.SetUpdateTime(time.Now().UTC())
	c.Metadata.SetState(false)
}

// Restore set active flag to true and sets UpdateTime metadata to current time in UTC
func (c *Category) Restore() {
	c.Metadata.SetUpdateTime(time.Now().UTC())
	c.Metadata.SetState(true)
}
