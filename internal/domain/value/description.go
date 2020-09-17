package value

import "github.com/alexandria-oss/common-go/exception"

type Description struct {
	value string
}

func NewDescription(d string) (*Description, error) {
	desc := new(Description)
	if err := desc.Set(d); err != nil {
		return nil, err
	}

	return desc, nil
}

func (d Description) Get() string {
	return d.value
}

func (d *Description) Set(description string) error {
	memo := d.value
	d.value = description

	if err := d.IsValid(); err != nil {
		d.value = memo
		return err
	}

	return nil
}

func (d Description) IsValid() error {
	// - Range from 1 to 512
	if d.value != "" && len(d.value) > 512 {
		return exception.NewFieldRange("description", "1", "512")
	}

	return nil
}
