package value

import (
	"strings"

	"github.com/alexandria-oss/common-go/exception"
)

// Color RGB-color scheme
type Color struct {
	value string
}

func NewColor(color string) (*Color, error) {
	t := new(Color)
	if err := t.Set(color); err != nil {
		return nil, err
	}

	return t, nil
}

func (t Color) Get() string {
	return t.value
}

func (t *Color) Set(color string) error {
	memoized := t.value
	t.value = strings.ToUpper(color)

	if err := t.IsValid(); err != nil {
		t.value = memoized
		return err
	}

	return nil
}

func (t Color) IsValid() error {
	isColor := t.value != "" && (t.value != "RED" && t.value != "BLUE" && t.value != "YELLOW" && t.value != "PINK" && t.value != "GREEN")
	if isColor {
		return exception.NewFieldFormat("color", "[Red, Blue, Yellow, Pink, Green)")
	}

	return nil
}
