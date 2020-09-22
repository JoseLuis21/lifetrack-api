package aggregate

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
)

type Category struct {
	root *entity.Category
}

func (c Category) IsValid() error {
	if err := c.root.IsValid(); err != nil {
		return err
	}

	return nil
}

func (c *Category) SetRoot(r *entity.Category) {
	c.root = r
}

func (c *Category) GetRoot() *entity.Category {
	return c.root
}
