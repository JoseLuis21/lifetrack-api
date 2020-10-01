package value

import "github.com/lucsky/cuid"

// CUID is a collision resistant unique identifier for cloud-native applications
type CUID struct {
	value string
}

func NewCUID() *CUID {
	c := new(CUID)
	c.Generate()
	return c
}

func (c *CUID) Generate() {
	c.value = cuid.New()
}

func (c *CUID) Set(id string) error {
	// cuid package does not contains a parser, cannot do anything to validate
	c.value = id
	return nil
}

func (c CUID) Get() string {
	return c.value
}
