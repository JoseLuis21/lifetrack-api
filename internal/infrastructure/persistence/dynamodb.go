package persistence

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewDynamoConn get a dynamodb connection from aws session pool
func NewDynamoConn(sess *session.Session, region string) *dynamodb.DynamoDB {
	return dynamodb.New(sess, &aws.Config{
		Region: aws.String(region),
	})
}
