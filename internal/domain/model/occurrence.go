package model

import (
	"encoding/json"

	"github.com/alexandria-oss/common-go/exception"
)

// Occurrence default-read model
type Occurrence struct {
	ID            string  `json:"id"`
	StartTime     int64   `json:"start_time"`
	EndTime       int64   `json:"end_time"`
	TotalDuration float64 `json:"total_duration"`
	Activity      string  `json:"activity"`
	CreateTime    int64   `json:"create_time"`
	UpdateTime    int64   `json:"update_time"`
	Active        bool    `json:"active"`
}

// MarshalBinary converts current model to binary data (JSON)
func (o Occurrence) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(o)
	if err != nil {
		return nil, exception.NewFieldFormat("occurrence", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to model
func (o *Occurrence) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, o); err != nil {
		return exception.NewFieldFormat("occurrence", "json")
	}

	return nil
}
