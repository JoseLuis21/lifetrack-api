package command

import (
	"context"
	"github.com/neutrinocorp/life-track-api/internal/application/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/event_factory"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// RemoveCategory request a category removal operation
type RemoveCategory struct {
	Ctx context.Context
	ID  string
}

// RemoveCategoryHandler handle category removal operations
type RemoveCategoryHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewRemoveCategoryHandler creates a new remove category command handler implementation
func NewRemoveCategoryHandler(r repository.Category, b event.Bus) *RemoveCategoryHandler {
	return &RemoveCategoryHandler{repo: r, bus: b}
}

func (h RemoveCategoryHandler) Invoke(cmd RemoveCategory) error {
	// Business ops
	id := value.CUID{}
	err := id.Set(cmd.ID)
	if err != nil {
		return err
	}

	// Get data
	c, err := h.repo.FetchByID(cmd.Ctx, id)
	if err != nil {
		return err
	}

	// If already deactivated, then skip
	if c.Active == false {
		return nil
	}

	// Parse primitive struct to domain aggregate
	category, err := adapter.CategoryAdapter{}.ToAggregate(*c)
	if err != nil {
		return err
	}
	// Store snapshot if rollback is needed
	snapshot := category

	// Update updateTime field
	category.Remove()

	// Persist changes
	err = h.repo.Replace(cmd.Ctx, *category)
	if err != nil {
		return err
	}

	// Publish events to message broker concurrent-safe
	errC := make(chan error)
	category.RecordEvent(event_factory.NewCategoryRemoved(*category.GetRoot().ID))
	go func() {
		if err = h.bus.Publish(cmd.Ctx, category.PullEvents()...); err != nil {
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
