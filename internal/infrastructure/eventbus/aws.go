package eventbus

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
)

// AWS is the event.Bus implementation using AWS SNS and SQS
type AWS struct {
	sess   *session.Session
	region string
	mu     *sync.Mutex
}

// NewAWS creates a concrete struct of AWS event bus
func NewAWS(s *session.Session, cfg infrastructure.Configuration) *AWS {
	return &AWS{
		sess:   s,
		region: cfg.EventBus.AWS.Region,
		mu:     new(sync.Mutex),
	}
}

func (b *AWS) Publish(ctx context.Context, e ...event.Domain) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(e) == 0 {
		return exception.NewRequiredField("domain event")
	}

	svc := NewSNSConn(b.sess, b.region)
	for _, ev := range e {
		ev.TopicToUnderscore()
		// Get topic Arn before publish
		topicArn, err := b.getTopicArn(ctx, svc, ev.Topic)
		if err != nil {
			return err
		}

		// Required SNS JSON struct
		defMsg := struct {
			Message string `json:"default"`
		}{
			Message: string(ev.Body),
		}

		eventJSON, err := json.Marshal(defMsg)
		if err != nil {
			return exception.NewFieldFormat("event_body", "binary or json")
		}

		_, err = svc.PublishWithContext(ctx, &sns.PublishInput{
			Message:           aws.String(string(eventJSON)),
			MessageAttributes: b.generateMessageAttr(ev),
			MessageStructure:  aws.String("json"),
			TopicArn:          aws.String(topicArn),
		})
		if err != nil {
			// fromSNSError already verifies nil err, but we cannot return nil from here,
			// otherwise for loop would be stopped if nil error was found
			return b.fromSNSError(err, ev)
		}
	}

	return nil
}

func (b *AWS) SubscribeTo(ctx context.Context, t event.Topic) (chan *event.Domain, error) {
	svc := NewSQSConn(b.sess, b.region)

	queueURL, err := b.getQueueURL(ctx, svc, string(t))
	if err != nil {
		return nil, err
	}

	// Long-polling strategy
	o, err := svc.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: aws.StringSlice([]string{
			"SentTimestamp",
		}),
		MaxNumberOfMessages: aws.Int64(1),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
		WaitTimeSeconds: aws.Int64(20),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sqs.ErrCodeOverLimit:
				return nil, exception.NewNetworkCall("aws sns topic "+string(t), b.region)
			}
		}

		return nil, err
	}

	// Use adapter func
	evChan := make(chan *event.Domain)
	e, err := b.getDomainEvent(o.Messages[0])
	if err != nil {
		return nil, err
	}
	evChan <- e

	return evChan, nil
}

// getTopicArn returns the given topic ARN from given session's AWS resources
func (b AWS) getTopicArn(ctx context.Context, svc *sns.SNS, topic string) (string, error) {
	nextToken := ""

	for {
		result, err := svc.ListTopicsWithContext(ctx, &sns.ListTopicsInput{
			NextToken: aws.String(nextToken),
		})
		if err != nil {
			return "", exception.NewNotFound("topics")
		}

		// Up to 100 topics
		for _, t := range result.Topics {
			// Search for given topic
			spl := strings.Split(*t.TopicArn, ":")
			if spl[len(spl)-1] == topic {
				return *t.TopicArn, nil
			}
		}

		// If no more to fetch, then break
		if result.NextToken == nil || *result.NextToken == "" {
			break
		}

		nextToken = *result.NextToken
	}

	return "", exception.NewNotFound("topic " + topic)
}

// getQueueURL returns a queue URL for the given queue name
func (b AWS) getQueueURL(ctx context.Context, svc *sqs.SQS, queue string) (string, error) {
	result, err := svc.GetQueueUrlWithContext(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	})
	if err != nil || result.QueueUrl == nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sqs.ErrCodeQueueDoesNotExist:
				return "", exception.NewNotFound("queue " + queue)
			case sqs.ErrCodeOverLimit:
				return "", exception.NewNetworkCall("aws sqs queue "+queue, b.region)
			}
		}

		return "", exception.NewNotFound("topic queue")
	}

	return *result.QueueUrl, nil
}

// getDomainEvent adapts sqs.Message into a domain event
func (b AWS) getDomainEvent(msg *sqs.Message) (*event.Domain, error) {
	if msg == nil || msg.Body == nil || msg.ReceiptHandle == nil {
		return nil, exception.NewRequiredField("message")
	}

	e := new(event.Domain)
	if err := e.UnmarshalBinary([]byte(*msg.Body)); err != nil {
		return nil, err
	}
	e.Acknowledge = *msg.ReceiptHandle

	return e, nil
}

// fromSNSError parses custom AWS errors to domain errors
func (b AWS) fromSNSError(err error, ev event.Domain) error {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sns.ErrCodeResourceNotFoundException:
				return exception.NewNotFound(ev.Topic)
			case sns.ErrCodeInvalidParameterException:
				return exception.NewFieldFormat(ev.Topic+" parameter", "valid topic parameter")
			case sns.ErrCodeInvalidParameterValueException:
				return exception.NewFieldFormat(ev.Topic+" parameter", "valid topic parameter value")
			case sns.ErrCodeThrottledException:
				return exception.NewNetworkCall("aws sns topic "+ev.Topic, b.region)
			}
		}

		return err
	}

	return nil
}

// generateMessageAttr set SNS message's attributes from event.Domain
func (b AWS) generateMessageAttr(ev event.Domain) map[string]*sns.MessageAttributeValue {
	attr := map[string]*sns.MessageAttributeValue{
		"id": {
			DataType:    aws.String("String"),
			StringValue: aws.String(ev.ID),
		},
		"service": {
			DataType:    aws.String("String"),
			StringValue: aws.String(ev.Service),
		},
		"action": {
			DataType:    aws.String("String"),
			StringValue: aws.String(ev.Action),
		},
		"aggregate_id": {
			DataType:    aws.String("String"),
			StringValue: aws.String(ev.AggregateID),
		},
		"aggregate_name": {
			DataType:    aws.String("String"),
			StringValue: aws.String(ev.AggregateName),
		},
		"publish_time": {
			DataType:    aws.String("Number"),
			StringValue: aws.String(strconv.FormatInt(ev.PublishTime.UTC().Unix(), 10)),
		},
	}

	if ev.Snapshot != nil {
		attr["snapshot"] = &sns.MessageAttributeValue{
			DataType:    aws.String("Binary"),
			BinaryValue: ev.Snapshot,
		}
	}

	return attr
}
