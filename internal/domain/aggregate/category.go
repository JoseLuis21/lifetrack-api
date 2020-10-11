package aggregate

import (
	"encoding/json"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
)

// Category is a group of n-activities tagged by a title
//	This is an aggregate root
type Category struct {
	category *entity.Category
	user     string
	events   []event.Domain
}

// Update mutates category data and sets UpdateTime metadata to current time in UTC
func (c *Category) Update(title, description, color string) error {
	if err := c.category.Update(title, description, color); err != nil {
		return err
	}

	return nil
}

// SetState sets the active flag to either false or true, this represents a soft-delete operation
func (c *Category) SetState(state bool) {
	if state {
		c.category.Restore()
		return
	}

	c.category.Remove()
}

// AssignUser attach a user to the current Category
func (c *Category) AssignUser(user string) error {
	if user == "" {
		return exception.NewRequiredField("category_user")
	}

	c.user = user
	return nil
}

// GetUser returns the current attached user
func (c Category) GetUser() string {
	return c.user
}

// ---- AGGREGATE ROOT IMPLEMENTATIONS ----

// IsValid validate all entities and value objects within aggregate bounds
func (c Category) IsValid() error {
	if err := c.category.IsValid(); err != nil {
		return err
	}

	return nil
}

// Set mutates the current entity.Category
func (c *Category) Set(category *entity.Category) {
	c.category = category
}

// Get returns the current entity.Category
func (c *Category) Get() *entity.Category {
	return c.category
}

// MarshalBinary converts current aggregate to binary data (JSON)
func (c Category) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(c)
	if err != nil {
		return nil, exception.NewFieldFormat("category_aggregate", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to aggregate
func (c *Category) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return exception.NewFieldFormat("category_aggregate", "json")
	}

	return nil
}

// RecordEvent records a new event.Domain into the aggregate
func (c *Category) RecordEvent(e event.Domain) {
	c.events = append(c.events, e)
}

// PullEvents returns all the event.Domain triggered in aggregate Category
func (c Category) PullEvents() []event.Domain {
	return c.events
}
