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

// ChangeCategoryState request a category state change
type ChangeCategoryState struct {
	Ctx   context.Context
	ID    string
	State bool
}

// ChangeCategoryStateHandler handles ChangeCategoryState commands
type ChangeCategoryStateHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewChangeCategoryStateHandler creates a new category state change command handler implementation
func NewChangeCategoryStateHandler(r repository.Category, b event.Bus) *ChangeCategoryStateHandler {
	return &ChangeCategoryStateHandler{
		repo: r,
		bus:  b,
	}
}

func (h ChangeCategoryStateHandler) Invoke(cmd ChangeCategoryState) error {
	// Business ops
	id := value.CUID{}
	if err := id.Set(cmd.ID); err != nil {
		return err
	}

	// Fetch data
	category, err := h.fetchCategory(cmd.Ctx, id)
	if err != nil {
		return err
	}
	snapshot := *category

	// Actual state change
	category.SetState(cmd.State)

	// Persist change
	err = h.repo.Replace(cmd.Ctx, *category)
	if err != nil {
		return err
	}

	// Publish events to message broker concurrent-safe
	h.setDomainEvent(category)
	return h.publishEvent(cmd.Ctx, category, snapshot)
}

func (h ChangeCategoryStateHandler) fetchCategory(ctx context.Context, id value.CUID) (*aggregate.Category, error) {
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

func (h ChangeCategoryStateHandler) setDomainEvent(ag *aggregate.Category) {
	if ag.Get().Metadata.GetState() {
		ag.RecordEvent(eventfactory.Category{}.NewCategoryRestored(*ag.Get().ID))
		return
	}

	ag.RecordEvent(eventfactory.Category{}.NewCategoryRemoved(*ag.Get().ID))
}

func (h ChangeCategoryStateHandler) publishEvent(ctx context.Context, ag *aggregate.Category, snapshot aggregate.Category) error {
	errC := make(chan error)
	go func() {
		if err := h.bus.Publish(ctx, ag.PullEvents()...); err != nil {
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
