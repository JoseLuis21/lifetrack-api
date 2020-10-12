package persistence

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/categorybuilder"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/readmodel"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

	db := NewDynamoConn(r.sess, r.tableRegion)
	id := util.GenerateDynamoID(r.schemaName, c.Get().ID.Get())
	_, err := db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(id),
			},
			"SK": {
				S: aws.String(id),
			},
			"title": {
				S: aws.String(c.Get().Title.Get()),
			},
			"description": {
				S: aws.String(c.Get().Description.Get()),
			},
			"theme": {
				S: aws.String(c.Get().Color.Get()),
			},
			"create_time": {
				N: aws.String(strconv.FormatInt(c.Get().Metadata.GetCreateTime().Unix(), 10)),
			},
			"update_time": {
				N: aws.String(strconv.FormatInt(c.Get().Metadata.GetUpdateTime().Unix(), 10)),
			},
			"active": {
				BOOL: aws.Bool(c.Get().Metadata.GetState()),
			},
			// DynamoDB GSI
			"GSIPK": {
				S: aws.String(util.GenerateDynamoID("User", c.GetUser())),
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

	db := NewDynamoConn(r.sess, r.tableRegion)
	res, err := db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:       r.generateKeyAttributes(id),
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

	var y categorybuilder.BuilderDynamo
	// Set strategy
	if criteria.User != "" {
		order := false
		if strings.ToUpper(criteria.OrderBy) == "ASC" {
			order = true
		}

		y = categorybuilder.NewUserDynamo(r.tableName, "GSIPK-index", r.schemaName).ByUser(criteria.User).
			Query(criteria.Query).Limit(limit).OrderBy(order).NextPage(token)
	} else {
		y = categorybuilder.NewCategoryDynamo(r.tableName, r.schemaName).
			Query(criteria.Query).Limit(limit).NextPage(token)
	}

	categories, nextPage, err := y.Do(ctx, NewDynamoConn(r.sess, r.tableRegion))
	if err != nil {
		return nil, "", r.getDomainError(err)
	} else if len(categories) == 0 {
		return nil, "", exception.NewNotFound("category")
	}

	return categories, nextPage, nil
}

func (r *CategoryDynamoRepository) Replace(ctx context.Context, c aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	db := NewDynamoConn(r.sess, r.tableRegion)
	_, err := db.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(c.Get().Title.Get()),
			},
			":d": {
				S: aws.String(c.Get().Description.Get()),
			},
			":th": {
				S: aws.String(c.Get().Color.Get()),
			},
			":u": {
				N: aws.String(strconv.FormatInt(c.Get().Metadata.GetUpdateTime().Unix(), 10)),
			},
			":a": {
				BOOL: aws.Bool(c.Get().Metadata.GetState()),
			},
		},
		Key:              r.generateKeyAttributes(*c.Get().ID),
		TableName:        aws.String(r.tableName),
		UpdateExpression: aws.String("SET title = :t, description = :d, update_time = :u, active = :a, theme = :th"),
	})

	return r.getDomainError(err)
}

func (r *CategoryDynamoRepository) HardRemove(ctx context.Context, id value.CUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	db := NewDynamoConn(r.sess, r.tableRegion)
	_, err := db.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key:       r.generateKeyAttributes(id),
		TableName: aws.String(r.tableName),
	})

	return r.getDomainError(err)
}

// generateKeyAttributes returns a primary and sort key map as a dynamo map
func (r CategoryDynamoRepository) generateKeyAttributes(id value.CUID) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String(util.GenerateDynamoID(r.schemaName, id.Get())),
		},
		"SK": {
			S: aws.String(util.GenerateDynamoID(r.schemaName, id.Get())),
		},
	}
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
