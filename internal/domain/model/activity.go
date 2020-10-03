package model

import (
	"encoding/json"
	"github.com/alexandria-oss/common-go/exception"
)

// Activity read model
type Activity struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Category   string `json:"category"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	Active     bool   `json:"active"`
}

// MarshalBinary converts current model to binary data (JSON)
func (a Activity) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(a)
	if err != nil {
		return nil, exception.NewFieldFormat("activity_aggregate", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to model
func (a *Activity) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, a); err != nil {
		return exception.NewFieldFormat("activity_aggregate", "json")
	}

	return nil
}
