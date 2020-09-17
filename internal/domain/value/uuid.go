package value

import "github.com/google/uuid"

type UUID struct {
	value uuid.UUID
}

// NewUUID generate a new UUID
func NewUUID() *UUID {
	id := new(UUID)
	id.Generate()
	return id
}

func (i UUID) Get() string {
	return i.value.String()
}

func (i *UUID) Set(id string) error {
	newID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	i.value = newID
	return nil
}

func (i *UUID) Generate() {
	i.value = uuid.New()
}
