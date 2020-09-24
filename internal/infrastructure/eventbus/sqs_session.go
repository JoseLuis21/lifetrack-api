package eventbus

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// NewSQSConn get an SQS connection from aws session pool
func NewSQSConn(sess *session.Session, region string) *sqs.SQS {
	return sqs.New(sess, &aws.Config{
		Region: aws.String(region),
	})
}
