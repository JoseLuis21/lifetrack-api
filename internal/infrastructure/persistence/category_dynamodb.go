package persistence

import (
	"context"
	"strconv"
	"sync"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/readmodel"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
)

// CategoryDynamoRepository is the repository.Category implementation using AWS DynamoDB
type CategoryDynamoRepository struct {
	sess        *session.Session
	tableName   string
	tableRegion string
	schemaName  string
	mu          *sync.RWMutex
}

func NewCategoryDynamoRepository(s *session.Session, cfg infrastructure.Configuration) *CategoryDynamoRepository {
	return &CategoryDynamoRepository{
		sess:        s,
		tableName:   cfg.Persistence.DynamoDB.Table,
		tableRegion: cfg.Persistence.DynamoDB.Region,
		schemaName:  "Category",
		mu:          new(sync.RWMutex),
	}
}

func (r CategoryDynamoRepository) Save(ctx context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	svc := NewDynamoConn(r.sess, r.tableRegion)
	id := util.GenerateDynamoID(r.schemaName, c.GetRoot().ID.Get())
	_, err := svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(id),
			},
			"SK": {
				S: aws.String(id),
			},
			"title": {
				S: aws.String(c.GetRoot().Title.Get()),
			},
			"description": {
				S: aws.String(c.GetRoot().Description.Get()),
			},
			"theme": {
				S: aws.String(c.GetRoot().Theme.Get()),
			},
			"create_time": {
				N: aws.String(strconv.FormatInt(c.GetRoot().Metadata.GetCreateTime().Unix(), 10)),
			},
			"update_time": {
				N: aws.String(strconv.FormatInt(c.GetRoot().Metadata.GetUpdateTime().Unix(), 10)),
			},
			"active": {
				BOOL: aws.Bool(c.GetRoot().Metadata.GetState()),
			},
			// DynamoDB GSI
			"GSIPK": {
				S: aws.String(util.GenerateDynamoID("User", c.GetRoot().User)),
			},
			"GSISK": {
				S: aws.String(id),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
		TableName:    aws.String(r.tableName),
	})

	return r.getDomainError(err)
}

func (r CategoryDynamoRepository) FetchByID(ctx context.Context, id value.CUID) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc := NewDynamoConn(r.sess, r.tableRegion)
	res, err := svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(util.GenerateDynamoID(r.schemaName, id.Get())),
			},
			"SK": {
				S: aws.String(util.GenerateDynamoID(r.schemaName, id.Get())),
			},
		},
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, r.getDomainError(err)
	}

	if res.Item == nil {
		return nil, exception.NewNotFound("category")
	}

	m := new(readmodel.CategoryDynamo)
	err = dynamodbattribute.UnmarshalMap(res.Item, m)
	if err != nil {
		return nil, err
	}

	return m.ToModel(), nil
}

func (r CategoryDynamoRepository) Fetch(ctx context.Context, token string, limit int64,
	criteria shared.CategoryCriteria) ([]*model.Category, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nextTokenMap map[string]*dynamodb.AttributeValue
	if token != "" {
		nextTokenMap = map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(token),
			},
			"SK": {
				S: aws.String(token),
			},
		}
	}

	exp, _ := r.buildFilter(criteria)

	svc := NewDynamoConn(r.sess, r.tableRegion)
	res, err := svc.ScanWithContext(ctx, &dynamodb.ScanInput{
		ConsistentRead:            aws.Bool(false),
		ExclusiveStartKey:         nextTokenMap,
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		Limit:                     aws.Int64(limit),
		TableName:                 aws.String(r.tableName),
	})
	if err != nil {
		return nil, "", r.getDomainError(err)
	}

	if len(res.Items) == 0 {
		return nil, "", exception.NewNotFound("category")
	}

	categories := make([]*model.Category, 0)

	for _, i := range res.Items {
		c := new(readmodel.CategoryDynamo)
		err = dynamodbattribute.UnmarshalMap(i, c)
		if err != nil {
			return nil, "", err
		}
		categories = append(categories, c.ToModel())
	}

	/*
		err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &categories)
		if err != nil {
			return nil, "", err
		}*/

	nextToken := ""
	if t := res.LastEvaluatedKey["PK"]; t != nil {
		nextToken = *t.S
	}

	return categories, nextToken, nil
}

func (r *CategoryDynamoRepository) Replace(ctx context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	svc := NewDynamoConn(r.sess, r.tableRegion)
	_, err := svc.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(c.GetRoot().Title.Get()),
			},
			":d": {
				S: aws.String(c.GetRoot().Description.Get()),
			},
			":th": {
				S: aws.String(c.GetRoot().Theme.Get()),
			},
			":u": {
				N: aws.String(strconv.FormatInt(c.GetRoot().Metadata.GetUpdateTime().Unix(), 10)),
			},
			":a": {
				BOOL: aws.Bool(c.GetRoot().Metadata.GetState()),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(util.GenerateDynamoID(r.schemaName, c.GetRoot().ID.Get())),
			},
			"SK": {
				S: aws.String(util.GenerateDynamoID(r.schemaName, c.GetRoot().ID.Get())),
			},
		},
		TableName:        aws.String(r.tableName),
		UpdateExpression: aws.String("SET title = :t, description = :d, update_time = :u, active = :a, theme = :th"),
	})

	return r.getDomainError(err)
}

func (r *CategoryDynamoRepository) HardRemove(ctx context.Context, id value.CUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	svc := NewDynamoConn(r.sess, r.tableRegion)
	_, err := svc.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(util.GenerateDynamoID(r.schemaName, id.Get())),
			},
			"SK": {
				S: aws.String(util.GenerateDynamoID(r.schemaName, id.Get())),
			},
		},
		TableName: aws.String(r.tableName),
	})

	return r.getDomainError(err)
}

// getDomainError returns a valid domain error from awserr dynamodb exceptions
func (r CategoryDynamoRepository) getDomainError(err error) error {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return exception.NewNotFound("category")
			case dynamodb.ErrCodeIndexNotFoundException:
				return exception.NewNotFound("category_id")
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return exception.NewFieldFormat("category_conditional", "valid query conditional field")
			case dynamodb.ErrCodeRequestLimitExceeded:
				return exception.NewNetworkCall("aws dynamodb table "+r.tableName, r.tableRegion)
			}
		}

		return err
	}

	return nil
}

// buildFilter constructs category fetch query criteria
func (r CategoryDynamoRepository) buildFilter(c shared.CategoryCriteria) (expression.Expression, error) {
	conditions := []expression.ConditionBuilder{expression.BeginsWith(expression.Name("PK"),
		r.schemaName)}

	if c.User != "" {
		conditions = append(conditions, expression.Equal(expression.Name("GSIPK"), expression.Value(
			util.GenerateDynamoID("User", c.User))))
	}
	if c.Title != "" {
		conditions = append(conditions, expression.Equal(expression.Name("title"), expression.Value(c.Title)))
	}
	if c.Query != "" {
		conditions = append(conditions, expression.Contains(expression.Name("title"), c.Query))
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
