package command

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/application/factory"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

// AddCategory request a new category creation
type AddCategory struct {
	Ctx         context.Context
	Title       string
	User        string
	Description string
}

// AddCategoryHandler handles AddCategory commands
type AddCategoryHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewAddCategoryHandler creates a new add category command handler implementation
func NewAddCategoryHandler(r repository.Category, b event.Bus) *AddCategoryHandler {
	return &AddCategoryHandler{repo: r, bus: b}
}

func (h AddCategoryHandler) Invoke(cmd AddCategory) error {
	// Business ops
	c, err := factory.NewCategory(cmd.Title, cmd.User, cmd.Description)
	if err != nil {
		return err
	}
	if err = c.IsValid(); err != nil {
		return err
	}

	// Infrastructure ops
	if err = h.repo.Save(cmd.Ctx, *c); err != nil {
		return err
	}

	// Publish events to message broker concurrent-safe
	errC := make(chan error)
	go func() {
		if err = h.bus.Publish(cmd.Ctx, c.PullEvents()...); err != nil {
			// Rollback
			if errRoll := h.repo.HardRemove(cmd.Ctx, *c.GetRoot().ID); errRoll != nil {
				errC <- errRoll
				return
			}

			errC <- err
			return
		}

		errC <- nil
	}()

	select {
	case err = <-errC:
		return err
	}
}
