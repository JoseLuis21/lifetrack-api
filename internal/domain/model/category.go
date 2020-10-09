package model

import (
	"encoding/json"

	"github.com/alexandria-oss/common-go/exception"
)

// Category read model
type Category struct {
	ID          string `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	User        string `json:"user"`
	Theme       string `json:"theme,omitempty"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	Active      bool   `json:"active"`
}

// MarshalBinary converts current model to binary data (JSON)
func (c Category) MarshalBinary() ([]byte, error) {
	cJSON, err := json.Marshal(c)
	if err != nil {
		return nil, exception.NewFieldFormat("category aggregate", "json")
	}

	return cJSON, nil
}

// UnmarshalBinary injects binary data to model
func (c *Category) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return exception.NewFieldFormat("category aggregate", "json")
	}

	return nil
}
