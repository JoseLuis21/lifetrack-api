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

// ChangeState request a category state change
type ChangeState struct {
	Ctx   context.Context
	ID    string
	State bool
}

// ChangeStateHandler handles ChangeState commands
type ChangeStateHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewChangeStateHandler creates a new ChangeState command handler implementation
func NewChangeStateHandler(r repository.Category, b event.Bus) *ChangeStateHandler {
	return &ChangeStateHandler{
		repo: r,
		bus:  b,
	}
}

func (h ChangeStateHandler) Invoke(cmd ChangeState) error {
	// Business ops
	id := value.CUID{}
	if err := id.Set(cmd.ID); err != nil {
		return err
	}

	// Fetch data
	category, err := h.fetch(cmd.Ctx, id)
	if err != nil {
		return err
	}

	return h.persist(cmd, category)
}

// fetch retrieve given category
func (h ChangeStateHandler) fetch(ctx context.Context, id value.CUID) (*aggregate.Category, error) {
	m, err := h.repo.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Adapt struct to domain aggregate
	category, err := adapter.CategoryAdapter{}.ToAggregate(*m)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// setDomainEvent retrieves either restored or removed category domain event depending on current aggregate state
func (h ChangeStateHandler) setDomainEvent(ag *aggregate.Category) {
	if ag.Get().Metadata.GetState() {
		ag.RecordEvent(eventfactory.Category{}.NewCategoryRestored(*ag.Get().ID))
		return
	}

	ag.RecordEvent(eventfactory.Category{}.NewCategoryRemoved(*ag.Get().ID))
}

// publishEvent publishes a domain event concurrently
func (h ChangeStateHandler) publishEvent(ctx context.Context, category *aggregate.Category, snapshot aggregate.Category) error {
	h.setDomainEvent(category)
	errC := make(chan error)
	go func() {
		if err := h.bus.Publish(ctx, category.PullEvents()...); err != nil {
			// Rollback
			if errR := h.repo.Replace(ctx, snapshot); errR != nil {
				errC <- errR
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
func (h ChangeStateHandler) persist(cmd ChangeState, category *aggregate.Category) error {
	// Store snapshot if rollback is needed
	snapshot := *category

	// Actual state change
	category.SetState(cmd.State)

	// Persist change
	if err := h.repo.Replace(cmd.Ctx, *category); err != nil {
		return err
	}

	// Publish events to message broker concurrent-safe
	return h.publishEvent(cmd.Ctx, category, snapshot)
}
