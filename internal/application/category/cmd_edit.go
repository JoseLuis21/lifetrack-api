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

// Edit request a category mutation
type Edit struct {
	Ctx         context.Context
	ID          string
	Title       string
	Description string
	Theme       string
	Image       string
}

// EditHandler handles Edit commands
type EditHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewEditHandler creates a new Edit command handler implementation
func NewEditHandler(r repository.Category, b event.Bus) *EditHandler {
	return &EditHandler{repo: r, bus: b}
}

func (h EditHandler) Invoke(cmd Edit) error {
	// Business ops
	id := value.CUID{}
	err := id.Set(cmd.ID)
	if err != nil {
		return err
	}

	// Get data
	category, err := h.fetch(cmd.Ctx, id)
	if err != nil {
		return err
	}

	return h.persist(category, cmd)
}

// fetch retrieve given category
func (h EditHandler) fetch(ctx context.Context, id value.CUID) (*aggregate.Category, error) {
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
func (h EditHandler) publishEvent(ctx context.Context, ag *aggregate.Category, snapshot aggregate.Category) error {
	errC := make(chan error)
	go func() {
		if err := h.bus.Publish(ctx, ag.PullEvents()...); err != nil {
			// Rollback
			if errRoll := h.repo.Replace(ctx, snapshot); errRoll != nil {
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
func (h EditHandler) persist(category *aggregate.Category, cmd Edit) error {
	// Store snapshot if rollback is needed
	snapshot := *category

	// Update fields
	if err := category.Update(cmd.Title, cmd.Description, cmd.Theme, cmd.Image); err != nil {
		return err
	}

	// Replace in persistence layer
	if err := h.repo.Replace(cmd.Ctx, *category); err != nil {
		return err
	}

	// Publish domain events to message broker concurrent-safe
	category.RecordEvent(eventfactory.Category{}.NewCategoryUpdated(*category, snapshot))
	return h.publishEvent(cmd.Ctx, category, snapshot)
}
