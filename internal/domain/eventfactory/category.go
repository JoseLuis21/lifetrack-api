package eventfactory

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
)

type Category struct{}

var categoryName = "category"

// Added triggered when c new aggregate.Category has been added to a Category
func (c Category) Added(category model.Category) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: categoryName,
		Action:        create,
		AggregateID:   category.ID,
		Body:          category,
	})
}

// Updated triggered when an existing aggregate.Category has been updated
func (c Category) Updated(category model.Category) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: categoryName,
		Action:        update,
		AggregateID:   category.ID,
		Body:          category,
	})
}

// Restored triggered when an existing aggregate.Category has been restored
func (c Category) Restored(activityID string) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: categoryName,
		Action:        restore,
		AggregateID:   activityID,
		Body:          nil,
	})
}

// Removed triggered when an existing aggregate.Category has been removed softly
func (c Category) Removed(activityID string) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: categoryName,
		Action:        remove,
		AggregateID:   activityID,
		Body:          nil,
	})
}

// HardRemoved triggered when an existing aggregate.Category has been removed permanently
func (c Category) HardRemoved(activityID string) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: categoryName,
		Action:        hardRemove,
		AggregateID:   activityID,
		Body:          nil,
	})
}
