package command

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/application/factory"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// EditCategory request a category mutation
type EditCategory struct {
	Ctx         context.Context
	ID          string
	Title       string
	Description string
}

// EditCategoryHandler handles EditCategory commands
type EditCategoryHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewEditCategoryHandler creates a new EditCategory command handler implementation
func NewEditCategoryHandler(r repository.Category, b event.Bus) *EditCategoryHandler {
	return &EditCategoryHandler{repo: r, bus: b}
}

func (h EditCategoryHandler) Invoke(cmd EditCategory) error {
	// Business ops
	id := value.UUID{}
	err := id.Set(cmd.ID)
	if err != nil {
		return err
	}

	// Get data
	c, err := h.repo.FetchByID(cmd.Ctx, id)
	if err != nil {
		return err
	}

	// Parse primitive struct to domain aggregate
	category, err := factory.NewCategoryFromModel(*c)
	if err != nil {
		return err
	}
	// Store snapshot if rollback is needed
	snapshot := category

	// Update fields and trigger event
	if err = category.Update(cmd.Title, cmd.Description); err != nil {
		return err
	}

	// Replace in persistence layer
	if err = h.repo.Replace(cmd.Ctx, *category); err != nil {
		return err
	}

	// Publish domain events to message broker concurrent-safe
	errC := make(chan error)
	go func() {
		if err := h.bus.Publish(cmd.Ctx, category.PullEvents()...); err != nil {
			// Rollback
			if errRoll := h.repo.Replace(cmd.Ctx, *snapshot); errRoll != nil {
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
