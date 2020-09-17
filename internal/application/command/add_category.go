package command

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/application/factory"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

type AddCategory struct {
	Ctx         context.Context
	Title       string
	User        string
	Description string
}

type AddCategoryHandler struct {
	repo repository.Category
}

func NewAddCategoryHandler(r repository.Category) *AddCategoryHandler {
	return &AddCategoryHandler{repo: r}
}

func (h AddCategoryHandler) Handle(cmd AddCategory) error {
	c, err := factory.NewCategory(cmd.Title, cmd.User, cmd.Description)
	if err != nil {
		return err
	}

	if err := c.IsValid(); err != nil {
		return err
	}

	return h.repo.Save(cmd.Ctx, c)
}
