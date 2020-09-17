package aggregate

import (
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"time"
)

type Category struct {
	ID          *value.UUID
	Title       *value.Title
	Description *value.Description
	User        string
	CreateTime  time.Time
	UpdateTime  time.Time
}

func (c Category) IsValid() error {
	// - Required: title, id, user

	if c.ID == nil || c.ID.Get() == "" {
		return exception.NewRequiredField("id")
	} else if c.Title == nil || c.Title.Get() == "" {
		return exception.NewRequiredField("title")
	} else if c.User == "" {
		return exception.NewRequiredField("user")
	}

	return nil
}
