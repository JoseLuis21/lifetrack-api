package value

import (
	"time"

	"github.com/alexandria-oss/common-go/exception"
)

// AppointedTime is time set by user for each activity in minutes
type AppointedTime struct {
	value time.Duration
}

func NewAppointedTime(appointedTime int) (*AppointedTime, error) {
	ap := &AppointedTime{
		value: time.Minute * time.Duration(appointedTime),
	}

	if err := ap.IsValid(); err != nil {
		return nil, err
	}

	return ap, nil
}

func (t AppointedTime) Get() time.Duration {
	return t.value
}

func (t *AppointedTime) Set(appointedTime int) error {
	memoized := t.value

	t.value = time.Duration(appointedTime) * time.Minute
	if err := t.IsValid(); err != nil {
		t.value = memoized
		return err
	}

	return nil
}

func (t AppointedTime) IsValid() error {
	// rules
	// - from 10 minutes up to 1 year (525,600 min. aprox.)
	if t.value.Minutes() < 10 || t.value.Minutes() > 525600 {
		return exception.NewFieldRange("appointed_time", "10 minutes", "1 year")
	}

	return nil
}
