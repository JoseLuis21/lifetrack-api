package value

import (
	"strings"

	"github.com/alexandria-oss/common-go/exception"
)

// Title is the main text header of an entity
type Title struct {
	value     string
	fieldName string
}

// NewTitle create a new title
func NewTitle(fieldName, title string) (*Title, error) {
	t := new(Title)
	t.SetFieldName(fieldName)
	if err := t.Set(title); err != nil {
		return nil, err
	}

	return t, nil
}

func (t Title) Get() string {
	return t.value
}

func (t *Title) Set(title string) error {
	memo := t.value
	t.value = strings.Title(title)

	if err := t.IsValid(); err != nil {
		t.value = memo
		return err
	}

	return nil
}

func (t Title) IsValid() error {
	// - Range from 1 to 256
	if t.value != "" && len(t.value) > 256 {
		return exception.NewFieldRange(t.fieldName, "1", "256")
	}

	return nil
}

func (t *Title) SetFieldName(field string) {
	if field != "" {
		t.fieldName = strings.ToLower(field)
	}

	t.fieldName = "title"
}
