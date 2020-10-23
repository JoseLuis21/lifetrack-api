package value

import (
	"strings"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
)

type Image struct {
	value string
}

func NewImage(image string) (*Image, error) {
	i := &Image{value: image}
	if err := i.IsValid(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i Image) Get() string {
	return i.value
}

func (i *Image) Set(image string) error {
	memoized := i.value
	i.value = image
	if err := i.IsValid(); err != nil {
		i.value = memoized
		return err
	}

	return nil
}

func (i Image) IsValid() error {
	// rules
	// 1.	valid domain (e.g. prefix = "https://cdn.damascus-engineering.com/lifetrack/static/")
	// 2.	valid image extension (e.g. suffix [.jpeg, .jpg, .webp)

	domain := shared.CDNDomain + "/lifetrack/static/"

	// boolean operators/expressions
	isDomain := i.value != "" && !strings.HasPrefix(i.value, domain)

	// state machine
	if isDomain {
		return exception.NewFieldFormat("image", "valid domain ("+domain+")")
	} else if i.isImage() {
		return exception.NewFieldFormat("image", "valid picture format [.jpeg, .jpg, .webp)")
	}

	return nil
}

func (i Image) isImage() bool {
	return i.value != "" && (!strings.HasSuffix(i.value, ".jpg") && !strings.HasSuffix(i.value, ".jpeg") &&
		!strings.HasSuffix(i.value, ".webp"))
}
