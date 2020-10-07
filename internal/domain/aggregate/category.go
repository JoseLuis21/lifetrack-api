package aggregate

import (
	"encoding/json"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
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

// Update mutates editable data and sets UpdateTime metadata to current time in UTC
func (c *Category) Update(title, description, theme string) error {
	if err := c.root.Update(title, description, theme); err != nil {
		return err
	}

	return nil
}

// Remove sets active flag to false and sets UpdateTime metadata to current time in UTC
func (c *Category) Remove() {
	c.root.Remove()
}

// Restore set active flag to true and sets UpdateTime metadata to current time in UTC
func (c *Category) Restore() {
	c.root.Restore()
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
	cJSON, err := json.Marshal(c)
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
func (c Category) PullEvents() []event.Domain {
	return c.events
}
