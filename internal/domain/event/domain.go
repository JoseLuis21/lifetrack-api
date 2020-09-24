package event

import (
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/google/uuid"
	"strings"
	"time"
)

// Topic represents a topic name, must be used to subscribe to events
type Topic string

// Domain represents a domain event, triggered when something happened inside each context
type Domain struct {
	ID            string    `json:"correlation_id"`
	Topic         string    `json:"topic"`
	Service       string    `json:"service"`
	Action        string    `json:"action"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateName string    `json:"aggregate_name"`
	Body          []byte    `json:"body"`
	Snapshot      []byte    `json:"snapshot"`
	PublishTime   time.Time `json:"publish_time"`
	Acknowledge   string    `json:"-"`
}

// DomainArgsDTO Data Transfer object required to create Domain events
type DomainArgsDTO struct {
	Service       string
	Action        string
	AggregateID   string
	AggregateName string
	Body          encoding.BinaryMarshaler
	Snapshot      encoding.BinaryMarshaler
}

// NewDomain generates a new Domain event
func NewDomain(args DomainArgsDTO) (*Domain, error) {
	if args.Service == "" {
		return nil, exception.NewRequiredField("service")
	} else if args.AggregateName == "" {
		return nil, exception.NewRequiredField("aggregate_name")
	} else if args.Action == "" {
		return nil, exception.NewRequiredField("action")
	} else if args.AggregateID == "" {
		return nil, exception.NewRequiredField("aggregate_id")
	}

	var bodyBinary []byte
	bodyBinary = nil
	if args.Body != nil {
		d, err := args.Body.MarshalBinary()
		if err != nil {
			return nil, exception.NewFieldFormat("body", "binary or json")
		}
		bodyBinary = d
	}

	var snapshotBinary []byte
	snapshotBinary = nil
	if args.Snapshot != nil {
		d, err := args.Snapshot.MarshalBinary()
		if err != nil {
			return nil, exception.NewFieldFormat("snapshot", "binary or json")
		}
		snapshotBinary = d
	}

	return &Domain{
		ID: uuid.New().String(),
		Topic: strings.ToLower(fmt.Sprintf("%s.%s.%s", args.Service, args.AggregateName,
			args.Action)),
		Service:       strings.ToLower(args.Service),
		Action:        strings.ToUpper(args.Action),
		AggregateID:   args.AggregateID,
		AggregateName: strings.ToLower(args.AggregateName),
		Body:          bodyBinary,
		Snapshot:      snapshotBinary,
		PublishTime:   time.Now(),
		Acknowledge:   "",
	}, nil
}

// TopicToUnderscore formats current Topic value to underscore format
//	e.g. foo.bar.created -> foo_bar_created
func (d *Domain) TopicToUnderscore() {
	d.Topic = strings.ToLower(fmt.Sprintf("%s_%s_%s", d.Service, d.AggregateName,
		d.Action))
}

// MarshalBinary converts current Domain event into bytes
func (d Domain) MarshalBinary() ([]byte, error) {
	dJSON, err := json.Marshal(d)
	if err != nil {
		return nil, exception.NewFieldFormat("domain event", "json")
	}

	return dJSON, nil
}

// UnmarshalBinary sets current Domain event from bytes
func (d *Domain) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return exception.NewFieldFormat("domain event", "json")
	}

	return nil
}
