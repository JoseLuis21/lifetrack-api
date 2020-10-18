package value

import (
	"errors"
	"testing"

	"github.com/alexandria-oss/common-go/exception"
)

func TestNewImage(t *testing.T) {
	// DOMAIN ERRORS
	// a.	invalid domain
	if _, err := NewImage("http://foo.com/bar/image"); !errors.Is(err, exception.FieldFormat) {
		t.Fatal("failed to create image, expected field format exception")
	}

	// b.	invalid image format
	if _, err := NewImage("https://cdn.damascus-engineering.com/lifetrack/static/image.png"); !errors.Is(err, exception.FieldFormat) {
		t.Fatal("failed to create image, expected field format exception")
	}

	// VALID SCENARIOS
	// a.	image with jpg extension
	if _, err := NewImage("https://cdn.damascus-engineering.com/lifetrack/static/image.jpg"); err != nil {
		t.Log(exception.GetDescription(err))
		t.Fatal("failed to create image, expected nil exception")
	}

	// b.	image with jpeg extension
	if _, err := NewImage("https://cdn.damascus-engineering.com/lifetrack/static/image.jpeg"); err != nil {
		t.Log(exception.GetDescription(err))
		t.Fatal("failed to create image, expected nil exception")
	}

	// c.	image with webp extension
	if _, err := NewImage("https://cdn.damascus-engineering.com/lifetrack/static/image.webp"); err != nil {
		t.Log(exception.GetDescription(err))
		t.Fatal("failed to create image, expected nil exception")
	}
}
