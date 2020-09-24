package eventbus

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// NewSNSConn get an SNS connection from aws session pool
func NewSNSConn(sess *session.Session, region string) *sns.SNS {
	return sns.New(sess, &aws.Config{
		Region: aws.String(region),
	})
}
