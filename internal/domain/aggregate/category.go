package aggregate

import (
	"encoding/json"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"time"
)

// Category activity container
type Category struct {
	root   *entity.Category
	events []event.Domain
}

// IsValid validate all entities and value objects within aggregate bounds
func (c Category) IsValid() error {
	if err := c.root.IsValid(); err != nil {
		return err
	}

	return nil
}

// Update triggers updated event.Domain, mutates editable data and sets UpdateTime metadata to current time in UTC
func (c *Category) Update(title, description string) error {
	if err := c.root.Title.Set(title); title != "" && err != nil {
		return err
	}
	if err := c.root.Description.Set(description); description != "" && err != nil {
		return err
	}

	c.root.UpdateTime = time.Now().UTC()
	c.events = append(c.events, event.NewCategoryUpdated(*c))
	return nil
}

// Remove triggers removed event.Domain, sets active flag to false and sets UpdateTime metadata to current time in UTC
func (c *Category) Remove() {
	c.root.UpdateTime = time.Now().UTC()
	c.root.Active = false
	c.events = append(c.events, event.NewCategoryRemoved(*c.root.ID))
}

// Restore triggers restored event.Domain, set active flag to true and sets UpdateTime metadata to current time in UTC
func (c *Category) Restore() {
	c.root.UpdateTime = time.Now().UTC()
	c.root.Active = true
	c.events = append(c.events, event.NewCategoryRestored(*c.root.ID))
}

// SetRoot mutates the current aggregate root (category)
func (c *Category) SetRoot(r *entity.Category) {
	c.root = r
}

// GetRoot returns the current aggregate root (category)
func (c *Category) GetRoot() *entity.Category {
	return c.root
}

// MarshalBinary converts current aggregate to binary data (JSON)
func (c Category) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(&model.Category{
		ID:          c.root.ID.Get(),
		Title:       c.root.Title.Get(),
		Description: c.root.Description.Get(),
		User:        c.root.User,
		CreateTime:  c.root.CreateTime.Unix(),
		UpdateTime:  c.root.UpdateTime.Unix(),
		Active:      c.root.Active,
	})
	if err != nil {
		return nil, exception.NewFieldFormat("category aggregate", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to aggregate
func (c *Category) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return exception.NewFieldFormat("category aggregate", "json")
	}

	return nil
}

// RecordEvent records a new event.Domain into the aggregate
func (c *Category) RecordEvent(e event.Domain) {
	c.events = append(c.events, e)
}

// PullEvents returns all the event.Domain triggered in aggregate Category
func (c *Category) PullEvents() []event.Domain {
	return c.events
}
