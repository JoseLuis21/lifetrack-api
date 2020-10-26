package dynamocategory

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/configuration"
)

// fetchByUser strategy when criteria contains a user ID
type fetchByUser struct {
	cfg configuration.Configuration
	db  *dynamodb.DynamoDB
}

func (r fetchByUser) Do(ctx context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	return nil, "", nil
}
