package aggregate

import (
	"encoding/json"

	"github.com/neutrinocorp/life-track-api/internal/domain/value"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
)

// Activity is a user's habit, resides inside of a category if desired
type Activity struct {
	activity *entity.Activity
	category *value.CUID
	events   []event.Domain
}

// Update mutates editable data and sets UpdateTime metadata to current time in UTC
func (a *Activity) Update(title string, appointedTime int) error {
	if err := a.activity.Update(title, appointedTime); err != nil {
		return err
	}

	return nil
}

// SetState sets the active flag to either false or true, this represents a soft-delete operation
func (a *Activity) SetState(state bool) {
	if state {
		a.activity.Restore()
		return
	}

	a.activity.Remove()
}

// AssignCategory attach a Category to the current Activity
func (a *Activity) AssignCategory(category string) error {
	if category == "" {
		return exception.NewRequiredField("activity_category")
	}

	id := &value.CUID{}
	if err := id.Set(category); err != nil {
		return err
	}
	a.category = id

	return nil
}

// GetCategory returns the current attached Category
func (a Activity) GetCategory() string {
	return a.category.Get()
}

// ---- AGGREGATE IMPLEMENTATIONS ----

// IsValid validate all entities and value objects within aggregate bounds
func (a Activity) IsValid() error {
	if err := a.activity.IsValid(); err != nil {
		return err
	}
	return nil
}

// Set mutates the current entity.Activity
func (a *Activity) Set(activity *entity.Activity) {
	a.activity = activity
}

// Get returns the current entity.Activity
func (a *Activity) Get() *entity.Activity {
	return a.activity
}

// MarshalBinary converts current aggregate to binary data (JSON)
func (a Activity) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(a)
	if err != nil {
		return nil, exception.NewFieldFormat("activity_aggregate", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to aggregate
func (a *Activity) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, a); err != nil {
		return exception.NewFieldFormat("activity_aggregate", "json")
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
