package category

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/adapter"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/eventfactory"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
)

// Remove request a category removal operation (hard remove)
type Remove struct {
	Ctx context.Context
	ID  string
}

// RemoveHandler handle category removal operations
type RemoveHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewRemoveHandler creates a new Remove command handler implementation
func NewRemoveHandler(r repository.Category, b event.Bus) *RemoveHandler {
	return &RemoveHandler{repo: r, bus: b}
}

func (h RemoveHandler) Invoke(cmd Remove) error {
	// Business ops
	id := value.CUID{}
	err := id.Set(cmd.ID)
	if err != nil {
		return err
	}

	// Get data
	snapshot, err := h.fetch(cmd.Ctx, id)
	if err != nil {
		return err
	}

	return h.persist(cmd, snapshot)
}

// fetch retrieve given category
func (h RemoveHandler) fetch(ctx context.Context, id value.CUID) (*aggregate.Category, error) {
	// Get data
	c, err := h.repo.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Parse primitive struct to domain aggregate
	category, err := adapter.CategoryAdapter{}.ToAggregate(*c)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// publishEvent publishes a domain event concurrently
func (h RemoveHandler) publishEvent(ctx context.Context, snapshot aggregate.Category) error {
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

// persist saves all recorded changes to ecosystem's persistence
func (h RemoveHandler) persist(cmd Remove, snapshot *aggregate.Category) error {
	// Persist changes
	if err := h.repo.HardRemove(cmd.Ctx, *snapshot.Get().ID); err != nil {
		return err
	}

	// Publish events to message broker concurrent-safe
	return h.publishEvent(cmd.Ctx, *snapshot)
}
