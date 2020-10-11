package command

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/eventfactory"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// RemoveCategory request a category removal operation (hard remove)
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

	// Get data (required by snapshot if rollback is needed)
	c, err := h.repo.FetchByID(cmd.Ctx, id)
	if err != nil {
		return err
	}

	// Parse primitive struct to domain aggregate
	// Store snapshot if rollback is needed
	snapshot, err := adapter.CategoryAdapter{}.ToAggregate(*c)
	if err != nil {
		return err
	}

	// Persist changes
	err = h.repo.HardRemove(cmd.Ctx, *snapshot.Get().ID)
	if err != nil {
		return err
	}

	// Publish events to message broker concurrent-safe
	return h.publishEvent(cmd.Ctx, *snapshot)
}

func (h RemoveCategoryHandler) publishEvent(ctx context.Context, snapshot aggregate.Category) error {
	errC := make(chan error)
	go func() {
		if err := h.bus.Publish(ctx, eventfactory.Category{}.NewCategoryHardRemoved(snapshot)); err != nil {
			// Rollback
			if errRoll := h.repo.Save(ctx, snapshot); errRoll != nil {
				errC <- errRoll
				return
			}

			errC <- err
			return
		}

		errC <- nil
	}()

	return <-errC
}
