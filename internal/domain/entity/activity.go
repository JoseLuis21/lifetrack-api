package entity

import (
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"time"
)

// Activity represents a task in a category
type Activity struct {
	ID       *value.CUID
	Title    *value.Title
	Category *value.CUID
	Metadata *value.Metadata
}

// NewActivity creates a new activity entity
func NewActivity(title, category string) (*Activity, error) {
	t, err := value.NewTitle("activity_title", title)
	if err != nil {
		return nil, err
	}

	c := new(value.CUID)
	err = c.Set(category)
	if err != nil {
		return nil, err
	}

	return &Activity{
		ID:       value.NewCUID(),
		Title:    t,
		Category: c,
		Metadata: value.NewMetadata(),
	}, nil
}

func (a Activity) IsValid() error {
	// Required: id, title, category
	if a.ID == nil {
		return exception.NewRequiredField("activity_id")
	} else if a.Title == nil {
		return exception.NewRequiredField("activity_title")
	} else if a.Category == nil {
		return exception.NewRequiredField("activity_category")
	}

	return nil
}

// Update mutates editable data and sets UpdateTime metadata to current time in UTC
func (a *Activity) Update(title string) error {
	if err := a.Title.Set(title); title != "" && err != nil {
		return err
	}
	_ = a.Metadata.SetUpdateTime(time.Now().UTC())

	if err := a.IsValid(); err != nil {
		return err
	}

	return nil
}

// Remove sets active flag to false and sets UpdateTime metadata to current time in UTC
func (a *Activity) Remove() {
	_ = a.Metadata.SetUpdateTime(time.Now().UTC())
	_ = a.Metadata.SetState(false)
}

// Restore set active flag to true and sets UpdateTime metadata to current time in UTC
func (a *Activity) Restore() {
	_ = a.Metadata.SetUpdateTime(time.Now().UTC())
	_ = a.Metadata.SetState(true)
}
