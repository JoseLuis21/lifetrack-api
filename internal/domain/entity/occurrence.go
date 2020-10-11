package entity

import (
	"time"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// TODO: Turn start and end times into value objects

// Occurrence is a task done inside an activity
type Occurrence struct {
	ID        *value.CUID
	StartTime time.Time
	EndTime   time.Time
	// Duration in minutes (m)
	TotalDuration float64
	Metadata      *value.Metadata
}

// NewOccurrence creates a new occurrence with required values
func NewOccurrence(startTime, endTime string) (*Occurrence, error) {
	o := &Occurrence{
		ID:       value.NewCUID(),
		Metadata: value.NewMetadata(),
	}
	s, e, err := o.parseTimes(startTime, endTime)
	if err != nil {
		return nil, err
	}
	o.StartTime = s
	o.EndTime = e

	if err := o.IsValid(); err != nil {
		return nil, err
	}

	o.TotalDuration = e.Sub(s).Minutes()
	return o, nil
}

// parseTimes converts given time in primitive to built-in time struct
func (o Occurrence) parseTimes(s, e string) (start, end time.Time, err error) {
	start, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return start, end, exception.NewFieldFormat("activity_start_time", "time with RFC3339 format ("+time.RFC3339+")")
	}

	end, err = time.Parse(time.RFC3339, e)
	if err != nil {
		return start, end, exception.NewFieldFormat("activity_end_time", "time with RFC3339 format ("+time.RFC3339+")")
	}

	return start.UTC(), end.UTC(), err
}

// IsValid validates current Occurrence values
func (o Occurrence) IsValid() error {
	// rules
	// required: id, start_time, end_time
	// times up to 1 year at max
	switch {
	case o.ID.Get() == "":
		return exception.NewRequiredField("occurrence_id")
	case o.StartTime.IsZero():
		return exception.NewRequiredField("occurrence_start_time")
	case o.EndTime.IsZero():
		return exception.NewRequiredField("occurrence_end_time")
	case o.EndTime.Before(o.StartTime):
		return exception.NewFieldRange("occurrence_end_time", "start_time", "end_time")
	case o.StartTime.After(time.Now().UTC().Add(time.Hour * 8760)):
		return exception.NewFieldRange("occurrence_start_time", "1 minute", "1 year")
	case o.EndTime.After(time.Now().UTC().Add(time.Hour * 8760)):
		return exception.NewFieldRange("occurrence_end_time", "1 minute", "1 year")
	}

	return nil
}

// Update mutates data and sets UpdateTime metadata to current time in UTC
func (o *Occurrence) Update(startTime, endTime string) error {
	s, e, err := o.parseTimes(startTime, endTime)
	if err != nil {
		return err
	}
	o.StartTime = s
	o.EndTime = e
	o.Metadata.SetUpdateTime(time.Now().UTC())

	if err := o.IsValid(); err != nil {
		return err
	}
	return nil
}

// Remove sets active flag to false and sets UpdateTime metadata to current time in UTC
func (o *Occurrence) Remove() {
	o.Metadata.SetUpdateTime(time.Now().UTC())
	o.Metadata.SetState(false)
}

// Restore set active flag to true and sets UpdateTime metadata to current time in UTC
func (o *Occurrence) Restore() {
	o.Metadata.SetUpdateTime(time.Now().UTC())
	o.Metadata.SetState(true)
}
