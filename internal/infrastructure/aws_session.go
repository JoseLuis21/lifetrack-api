package infrastructure

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"sync"
)

var awsSingleton = new(sync.Once)
var awsSession *session.Session

func NewSession() *session.Session {
	if awsSession == nil {
		awsSingleton.Do(func() {
			awsSession = session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
		})
	}

	return awsSession
}
