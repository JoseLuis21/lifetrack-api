package category

import (
	"context"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// BuilderDynamo constructs and executes an AWS DynamoDB query using the
// given strategy (default, by user)
type BuilderDynamo interface {
	Do(ctx context.Context, db *dynamodb.DynamoDB) ([]*aggregate.Category, string, error)
}
