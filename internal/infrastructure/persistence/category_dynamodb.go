package persistence

import (
	"context"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/neutrinocorp/life-track-api/internal/domain"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"strconv"
	"strings"
	"sync"
)

type CategoryDynamoRepository struct {
	sess *session.Session
	cfg  infrastructure.Configuration
	mu   *sync.RWMutex
}

func NewCategoryDynamoRepository(s *session.Session, cfg infrastructure.Configuration) *CategoryDynamoRepository {
	return &CategoryDynamoRepository{
		sess: s,
		cfg:  cfg,
		mu:   new(sync.RWMutex),
	}
}

func (r CategoryDynamoRepository) Save(ctx context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// TODO: Use high-cardinality fields as keys since Dynamo is a key:value storage,
	// we store them within the same table or in new tables
	// we must do this using Dynamo Transactions
	r.mu.Unlock()
	exists, err := r.Exists(ctx, *c.GetRoot().Title, c.GetRoot().User)
	if err != nil {
		r.mu.Lock()
		return err
	} else if exists {
		r.mu.Lock()
		return exception.NewAlreadyExists("category")
	}
	r.mu.Lock()

	svc := NewDynamoConn(r.sess, r.cfg.Table.Region)
	_, err = svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"category_id": {
				S: aws.String(c.GetRoot().ID.Get()),
			},
			"title": {
				S: aws.String(c.GetRoot().Title.Get()),
			},
			"description": {
				S: aws.String(c.GetRoot().Description.Get()),
			},
			"user": {
				S: aws.String(c.GetRoot().User),
			},
			"create_time": {
				N: aws.String(strconv.FormatInt(c.GetRoot().CreateTime.Unix(), 10)),
			},
			"update_time": {
				N: aws.String(strconv.FormatInt(c.GetRoot().UpdateTime.Unix(), 10)),
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

func (r CategoryDynamoRepository) FetchByID(ctx context.Context, id value.UUID) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc := NewDynamoConn(r.sess, r.cfg.Table.Region)
	res, err := svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"category_id": {
				S: aws.String(id.Get()),
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

	m := new(model.Category)
	err = dynamodbattribute.UnmarshalMap(res.Item, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r CategoryDynamoRepository) Fetch(ctx context.Context, token string, limit int64, filter domain.FilterMap) ([]*model.Category, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nextTokenMap map[string]*dynamodb.AttributeValue
	nextTokenMap = nil
	if token != "" {
		nextTokenMap = map[string]*dynamodb.AttributeValue{
			"category_id": {
				S: aws.String(token),
			},
		}
	}

	exp, _ := r.buildFilter(filter)

	svc := NewDynamoConn(r.sess, r.cfg.Table.Region)
	res, err := svc.ScanWithContext(ctx, &dynamodb.ScanInput{
		AttributesToGet:           nil,
		ConditionalOperator:       nil,
		ConsistentRead:            aws.Bool(false),
		ExclusiveStartKey:         nextTokenMap,
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		IndexName:                 nil,
		Limit:                     aws.Int64(limit),
		ProjectionExpression:      nil,
		ReturnConsumedCapacity:    nil,
		ScanFilter:                nil,
		Segment:                   nil,
		Select:                    nil,
		TableName:                 aws.String(r.cfg.Table.Name),
		TotalSegments:             nil,
	})
	if err != nil {
		return nil, "", err
	}

	if len(res.Items) == 0 {
		return nil, "", exception.NewNotFound("category")
	}

	categories := make([]*model.Category, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &categories)
	if err != nil {
		return nil, "", err
	}

	nextToken := ""
	if t := res.LastEvaluatedKey["category_id"]; t != nil {
		nextToken = *t.S
	}

	return categories, nextToken, nil
}

func (r CategoryDynamoRepository) buildFilter(filter domain.FilterMap) (expression.Expression, error) {
	conditions := make([]expression.ConditionBuilder, 0)
	for k, v := range filter {
		switch {
		case k == "user" && v != "":
			conditions = append(conditions, expression.Equal(expression.Name("user"), expression.Value(v)))
			continue
		case k == "title" && v != "":
			conditions = append(conditions, expression.Equal(expression.Name("title"), expression.Value(v)))
			continue
		case k == "query" && v != "":
			conditions = append(conditions, expression.Contains(expression.Name("title"), v))
			continue
		}
	}

	if cLength := len(conditions); cLength >= 1 {
		if cLength == 2 {
			return expression.NewBuilder().WithFilter(expression.And(conditions[0], conditions[1])).Build()
		} else if cLength >= 3 {
			return expression.NewBuilder().WithFilter(expression.And(conditions[0], conditions[1], conditions[2:]...)).Build()
		}

		return expression.NewBuilder().WithFilter(conditions[0]).Build()
	}

	return expression.NewBuilder().Build()
}

func (r CategoryDynamoRepository) Exists(ctx context.Context, title value.Title, user string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc := NewDynamoConn(r.sess, r.cfg.Table.Region)

	exp, err := expression.NewBuilder().
		WithFilter(
			expression.And(expression.Name("title").Equal(expression.Value(strings.Title(title.Get()))),
				expression.Name("user").Equal(expression.Value(user)))).
		WithProjection(expression.NamesList(expression.Name("category_id"))).Build()
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
		Limit:                     nil,
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
	} else if res != nil && len(res.Items) >= 1 {
		return true, nil
	}

	return false, nil
}
