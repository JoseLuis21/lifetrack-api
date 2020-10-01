package value

import (
	"github.com/alexandria-oss/common-go/exception"
	"strings"
)

// Theme RGB-color scheme
type Theme struct {
	value string
}

func NewTheme(theme string) (*Theme, error) {
	t := new(Theme)
	if err := t.Set(theme); err != nil {
		return nil, err
	}

	return t, nil
}

func (t Theme) Get() string {
	return t.value
}

func (t *Theme) Set(theme string) error {
	memoized := t.value
	t.value = strings.ToUpper(theme)

	if err := t.IsValid(); err != nil {
		t.value = memoized
		return err
	}

	return nil
}

func (t Theme) IsValid() error {
	if t.value != "" && (t.value != "RED" && t.value != "BLUE" && t.value != "YELLOW" && t.value != "PINK" && t.value != "GREEN") {
		return exception.NewFieldFormat("theme", "[Red, Blue, Yellow, Pink, Green]")
	}

	return nil
}
