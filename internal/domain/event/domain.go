package event

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/google/uuid"
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

// DomainArgs Data Transfer object required to create Domain events
type DomainArgs struct {
	Service       string
	Action        string
	AggregateID   string
	AggregateName string
	Body          encoding.BinaryMarshaler
	Snapshot      encoding.BinaryMarshaler
}

// NewDomain generates a new Domain event
func NewDomain(args DomainArgs) (*Domain, error) {
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
	bin, err := encodeBinary("body", args.Body)
	if err != nil {
		return nil, err
	}
	bodyBinary = bin

	var snapshotBinary []byte
	bin, err = encodeBinary("snapshot", args.Snapshot)
	if err != nil {
		return nil, err
	}
	snapshotBinary = bin

	return &Domain{
		ID: uuid.New().String(),
		Topic: strings.ToLower(fmt.Sprintf("lt.%s.%s", args.AggregateName,
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

// encodeBinary parses a binary marshaller to binary array
func encodeBinary(field string, args encoding.BinaryMarshaler) ([]byte, error) {
	if args != nil {
		d, err := args.MarshalBinary()
		if err != nil {
			return nil, exception.NewFieldFormat(field, "binary or json")
		}
		return d, nil
	}

	return nil, nil
}

// TopicToUnderscore formats current Topic value to underscore format
//	e.g. lt.foo.created -> lt_foo_created
func (d *Domain) TopicToUnderscore() {
	d.Topic = strings.ToLower(fmt.Sprintf("lt_%s_%s", d.AggregateName,
		d.Action))
}

// TopicToDash formats current Topic value to dashed format
//	e.g. lt.foo.created -> lt-foo-created
func (d *Domain) TopicToDash() {
	d.Topic = strings.ToLower(fmt.Sprintf("lt-%s-%s", d.AggregateName,
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
