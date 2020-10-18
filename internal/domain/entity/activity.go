package entity

import (
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Activity is a user's habit, resides inside of a category if desired
type Activity struct {
	ID            *value.CUID
	Name          *value.Title
	AppointedTime *value.AppointedTime
	Metadata      *value.Metadata
}

// NewActivity creates a new activity entity
func NewActivity(name string, appointedTime int) (*Activity, error) {
	t, err := value.NewTitle("activity_name", name)
	if err != nil {
		return nil, err
	}

	act := &Activity{
		ID:       value.NewCUID(),
		Name:     t,
		Metadata: value.NewMetadata(),
	}

	apTime, err := value.NewAppointedTime(appointedTime)
	if err != nil {
		return nil, err
	}
	act.AppointedTime = apTime

	if err := act.IsValid(); err != nil {
		return nil, err
	}

	return act, nil
}

func (a Activity) IsValid() error {
	// a.	required fields [id, name)
	if a.ID == nil {
		return exception.NewRequiredField("activity_id")
	} else if a.Name == nil {
		return exception.NewRequiredField("activity_title")
	}

	return nil
}

// Update mutates data atomically and sets UpdateTime metadata to current time in UTC
func (a *Activity) Update(title string, appointedTime int) error {
	if err := a.Name.Set(title); title != "" && err != nil {
		return err
	}
	if appointedTime > 0 {
		apTime, err := value.NewAppointedTime(appointedTime)
		if err != nil {
			return err
		}
		a.AppointedTime = apTime
	}
	a.Metadata.SetUpdateTime(time.Now().UTC())

	if err := a.IsValid(); err != nil {
		return err
	}
	return nil
}

// Remove sets active flag to false and sets UpdateTime metadata to current time in UTC
func (a *Activity) Remove() {
	a.Metadata.SetUpdateTime(time.Now().UTC())
	a.Metadata.SetState(false)
}

// Restore set active flag to true and sets UpdateTime metadata to current time in UTC
func (a *Activity) Restore() {
	a.Metadata.SetUpdateTime(time.Now().UTC())
	a.Metadata.SetState(true)
}
