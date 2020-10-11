package aggregate

import (
	"encoding/json"

	"github.com/neutrinocorp/life-track-api/internal/domain/value"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
)

// Occurrence is a user's habit, resides inside of a category if desired
type Occurrence struct {
	occurrence *entity.Occurrence
	activity   *value.CUID
	events     []event.Domain
}

// Update mutates editable data and sets UpdateTime metadata to current time in UTC
func (o *Occurrence) Update(startTime, endTime string) error {
	if err := o.occurrence.Update(startTime, endTime); err != nil {
		return err
	}

	return nil
}

// SetState sets the active flag to either false or true, this represents a soft-delete operation
func (o *Occurrence) SetState(state bool) {
	if state {
		o.occurrence.Restore()
		return
	}

	o.occurrence.Remove()
}

// AssignActivity attach an Activity to the current Occurrence
func (o *Occurrence) AssignActivity(activity string) error {
	if activity == "" {
		return exception.NewRequiredField("occurrence_activity")
	}

	id := &value.CUID{}
	if err := id.Set(activity); err != nil {
		return err
	}
	o.activity = id

	return nil
}

// GetActivity returns the current attached Activity
func (o Occurrence) GetActivity() string {
	return o.activity.Get()
}

// ---- AGGREGATE IMPLEMENTATIONS ----

// IsValid validate all entities and value objects within aggregate bounds
func (o Occurrence) IsValid() error {
	if err := o.occurrence.IsValid(); err != nil {
		return err
	}
	return nil
}

// Set mutates the current entity.Occurrence
func (o *Occurrence) Set(Occurrence *entity.Occurrence) {
	o.occurrence = Occurrence
}

// Get returns the current entity.Occurrence
func (o *Occurrence) Get() *entity.Occurrence {
	return o.occurrence
}

// MarshalBinary converts current aggregate to binary data (JSON)
func (o Occurrence) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(o)
	if err != nil {
		return nil, exception.NewFieldFormat("occurrence_aggregate", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to aggregate
func (o *Occurrence) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, o); err != nil {
		return exception.NewFieldFormat("occurrence_aggregate", "json")
	}

	return nil
}

// RecordEvent records a new event.Domain into the aggregate
func (o *Occurrence) RecordEvent(e event.Domain) {
	o.events = append(o.events, e)
}

// PullEvents returns all the event.Domain triggered in aggregate Occurrence
func (o Occurrence) PullEvents() []event.Domain {
	return o.events
}
