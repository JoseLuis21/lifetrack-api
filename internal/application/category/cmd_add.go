package category

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/factory"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

// Add request a new category creation
type Add struct {
	Ctx         context.Context
	Title       string
	User        string
	Description string
	Theme       string
}

// AddHandler handles Add commands
type AddHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewAddHandler creates a new Add command handler implementation
func NewAddHandler(r repository.Category, b event.Bus) *AddHandler {
	return &AddHandler{repo: r, bus: b}
}

func (h AddHandler) Invoke(cmd Add) error {
	// Business ops
	c, err := factory.NewCategory(cmd.Title, cmd.User, cmd.Description, cmd.Theme)
	if err != nil {
		return err
	}

	// Infrastructure ops
	if err = h.repo.Save(cmd.Ctx, *c); err != nil {
		return err
	}

	// Publish events to message broker concurrent-safe
	return h.publishEvent(cmd.Ctx, c)
}

func (h AddHandler) publishEvent(ctx context.Context, c *aggregate.Category) error {
	errC := make(chan error)
	go func() {
		if err := h.bus.Publish(ctx, c.PullEvents()...); err != nil {
			// Rollback
			if errRoll := h.repo.HardRemove(ctx, *c.Get().ID); errRoll != nil {
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
