package persistence

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
)

// NewDynamoConn get a dynamodb connection from aws session pool
func NewDynamoConn(region string) *dynamodb.DynamoDB {
	return dynamodb.New(infrastructure.NewSession(), &aws.Config{
		Region: aws.String(region),
	})
}
