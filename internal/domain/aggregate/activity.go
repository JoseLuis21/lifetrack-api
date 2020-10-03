package aggregate

import (
	"encoding/json"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
)

// Activity contains activity behavior
type Activity struct {
	root   *entity.Activity
	events []event.Domain
}

// IsValid validate all entities and value objects within aggregate bounds
func (a Activity) IsValid() error {
	if err := a.root.IsValid(); err != nil {
		return err
	}

	return nil
}

// Update mutates editable data and sets UpdateTime metadata to current time in UTC
func (a *Activity) Update(title string) error {
	if err := a.root.Update(title); err != nil {
		return err
	}

	return nil
}

// Remove sets active flag to false and sets UpdateTime metadata to current time in UTC
func (a *Activity) Remove() {
	a.root.Remove()
}

// Restore set active flag to true and sets UpdateTime metadata to current time in UTC
func (c *Activity) Restore() {
	c.root.Restore()
}

// SetRoot mutates the current aggregate root (Activity)
func (a *Activity) SetRoot(r *entity.Activity) {
	a.root = r
}

// GetRoot returns the current aggregate root (Activity)
func (a *Activity) GetRoot() *entity.Activity {
	return a.root
}

// MarshalBinary converts current aggregate to binary data (JSON)
func (a Activity) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(a)
	if err != nil {
		return nil, exception.NewFieldFormat("activity aggregate", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to aggregate
func (a *Activity) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, a); err != nil {
		return exception.NewFieldFormat("activity aggregate", "json")
	}

	return nil
}

// RecordEvent records a new event.Domain into the aggregate
func (a *Activity) RecordEvent(e event.Domain) {
	a.events = append(a.events, e)
}

// PullEvents returns all the event.Domain triggered in aggregate Activity
func (a Activity) PullEvents() []event.Domain {
	return a.events
}
