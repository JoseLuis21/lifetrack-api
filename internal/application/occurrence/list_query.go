package occurrence

import (
	"context"

	"github.com/neutrinocorp/lifetrack-api/internal/domain/adapter"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/model"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
)

// ListQuery requests a set of model.Occurrence
type ListQuery struct {
	repo repository.Occurrence
}

// NewListQuery creates a ListQuery
func NewListQuery(r repository.Occurrence) *ListQuery {
	return &ListQuery{repo: r}
}

// Filter sets Occurrence fetching strategy
//	anti-corruption struct
type Filter struct {
	ActivityID string `json:"activity_id"`
	Limit      int64  `json:"limit"`
	Token      string `json:"token"`
}

func (q ListQuery) Query(ctx context.Context, filter Filter) ([]*model.Occurrence, string, error) {
	if filter.Limit == 0 {
		filter.Limit = 100
	}

	ocs, nextPage, err := q.repo.Fetch(ctx, repository.OccurrenceCriteria{
		Activity: filter.ActivityID,
		Limit:    filter.Limit,
		Token:    filter.Token,
	})
	return adapter.BulkUnmarshalPrimitiveOccurrence(ocs), nextPage, err
}
