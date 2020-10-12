package categorybuilder

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
)

// BuilderDynamo constructs and executes an AWS DynamoDB query using the
// given strategy (default, by user)
type BuilderDynamo interface {
	Do(ctx context.Context, db *dynamodb.DynamoDB) ([]*model.Category, string, error)
}
