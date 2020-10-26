package occurrence

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
)

// GetQuery requests a get Occurrence query
type GetQuery struct {
	repo repository.Occurrence
}

// NewGetQuery creates a new Get query
func NewGetQuery(r repository.Occurrence) *GetQuery {
	return &GetQuery{repo: r}
}

// Query handles Get requests
func (q GetQuery) Query(ctx context.Context, id string) (*model.Occurrence, error) {
	oc, _, err := q.repo.Fetch(ctx, repository.OccurrenceCriteria{ID: id})
	if err != nil {
		return nil, err
	}

	return oc[0].MarshalPrimitive(), nil
}
