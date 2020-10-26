package dbpool

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/remote"
)

// NewDynamoDB creates a new AWS DynamoDB connection pool
func NewDynamoDB(cfg configuration.Configuration) *dynamodb.DynamoDB {
	return dynamodb.New(remote.NewAWSSession(cfg.DynamoTable.Region))
}
