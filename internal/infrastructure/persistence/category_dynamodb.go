package persistence

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"strconv"
	"sync"
)

type CategoryDynamoRepository struct {
	sess *session.Session
	cfg  infrastructure.Configuration
	mu   *sync.RWMutex
}

func NewCategoryDynamoRepository(s *session.Session, cfg *infrastructure.Configuration) *CategoryDynamoRepository {
	return &CategoryDynamoRepository{
		sess: s,
		cfg:  *cfg,
		mu:   new(sync.RWMutex),
	}
}

func (r CategoryDynamoRepository) Save(ctx context.Context, c *aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.mu.Unlock()
	if exists, _ := r.Exists(ctx, c.Title.Get(), c.User); exists {
		r.mu.Lock()
		return exception.NewAlreadyExists("category")
	}
	r.mu.Lock()

	svc := NewDynamoConn(r.sess, r.cfg.Table.Region)
	_, err := svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"category_id": {
				S: aws.String(c.ID.Get()),
			},
			"title": {
				S: aws.String(c.Title.Get()),
			},
			"description": {
				S: aws.String(c.Description.Get()),
			},
			"user": {
				S: aws.String(c.User),
			},
			"create_time": {
				N: aws.String(strconv.FormatInt(c.CreateTime.Unix(), 10)),
			},
			"update_time": {
				N: aws.String(strconv.FormatInt(c.UpdateTime.Unix(), 10)),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
		TableName:    aws.String(r.cfg.Table.Name),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryDynamoRepository) FetchByID(ctx context.Context, id string) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc := NewDynamoConn(r.sess, r.cfg.Table.Region)

	res, err := svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"category_id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(r.cfg.Table.Name),
	})
	if err != nil {
		return nil, err
	}

	if res.Item == nil {
		return nil, exception.NewNotFound("category")
	}

	m := &model.Category{}
	err = dynamodbattribute.UnmarshalMap(res.Item, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r CategoryDynamoRepository) Exists(ctx context.Context, title string, user string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc := NewDynamoConn(r.sess, r.cfg.Table.Region)

	exp, err := expression.NewBuilder().WithFilter(
		expression.And(expression.Name("title").Equal(expression.Value(title)),
			expression.Name("user").Equal(expression.Value(user))),
	).WithProjection(expression.NamesList(expression.Name("category_id"))).Build()
	if err != nil {
		return false, err
	}

	res, err := svc.ScanWithContext(ctx, &dynamodb.ScanInput{
		AttributesToGet:           nil,
		ConditionalOperator:       nil,
		ConsistentRead:            aws.Bool(true),
		ExclusiveStartKey:         nil,
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		IndexName:                 nil,
		Limit:                     aws.Int64(1),
		ProjectionExpression:      exp.Projection(),
		ReturnConsumedCapacity:    nil,
		ScanFilter:                nil,
		Segment:                   nil,
		Select:                    nil,
		TableName:                 aws.String(r.cfg.Table.Name),
		TotalSegments:             nil,
	})
	if err != nil {
		return false, err
	}

	if res != nil && *res.Count >= 1 {
		return true, nil
	}

	return false, nil
}
