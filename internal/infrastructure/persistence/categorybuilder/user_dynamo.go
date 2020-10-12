package categorybuilder

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/readmodel"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"
)

// UserDynamo constructs and executes an AWS DynamoDB query fetching user's categories only
//	This is a concrete CategoryDynamo strategy implementation
type UserDynamo struct {
	user   string
	schema string
	exp    expression.Builder
	input  *dynamodb.QueryInput
}

func NewUserDynamo(tableName, indexName, schema string) *UserDynamo {
	return &UserDynamo{
		user:   "",
		schema: schema,
		exp:    expression.NewBuilder(),
		input: &dynamodb.QueryInput{
			IndexName: aws.String(indexName),
			Limit:     aws.Int64(100),
			TableName: aws.String(tableName),
		},
	}
}

func (b *UserDynamo) GetInput() *dynamodb.QueryInput {
	exp, _ := b.exp.Build()
	b.input.SetExpressionAttributeNames(exp.Names())
	b.input.SetExpressionAttributeValues(exp.Values())
	if exp.Filter() != nil {
		b.input.SetFilterExpression(*exp.Filter())
	}
	if exp.KeyCondition() != nil {
		b.input.SetKeyConditionExpression(*exp.KeyCondition())
	}

	return b.input
}

func (b *UserDynamo) ByUser(user string) *UserDynamo {
	if user != "" {
		b.user = user
		b.exp = b.exp.WithKeyCondition(expression.KeyAnd(expression.Key("GSIPK").
			Equal(expression.Value(util.GenerateDynamoID("User", user))),
			expression.KeyBeginsWith(expression.Key("GSISK"), b.schema)))
	}

	return b
}

func (b *UserDynamo) Query(keyword string) *UserDynamo {
	if keyword != "" {
		b.exp = b.exp.WithFilter(expression.Contains(expression.Name("title"), keyword))
	}

	return b
}

func (b *UserDynamo) Limit(l int64) *UserDynamo {
	if l > 0 {
		b.input.SetLimit(l)
	}

	return b
}

// OrderBy true -> asc, false -> desc
func (b *UserDynamo) OrderBy(o bool) *UserDynamo {
	b.input.SetScanIndexForward(o)
	return b
}

func (b *UserDynamo) NextPage(token string) *UserDynamo {
	if token != "" {
		b.input.SetExclusiveStartKey(map[string]*dynamodb.AttributeValue{
			"GSIPK": {
				S: aws.String(util.GenerateDynamoID("User", b.user)),
			},
			"GSISK": {
				S: aws.String(util.GenerateDynamoID(b.schema, token)),
			},
			"PK": {
				S: aws.String(util.GenerateDynamoID(b.schema, token)),
			},
			"SK": {
				S: aws.String(util.GenerateDynamoID(b.schema, token)),
			},
		})
	}

	return b
}

func (b UserDynamo) Do(ctx context.Context, db *dynamodb.DynamoDB) ([]*model.Category, string, error) {
	o, err := db.QueryWithContext(ctx, b.GetInput())
	if err != nil {
		return nil, "", err
	}

	categories := make([]*model.Category, 0)
	for _, i := range o.Items {
		c := new(readmodel.CategoryDynamo)
		err = dynamodbattribute.UnmarshalMap(i, c)
		if err != nil {
			return nil, "", err
		}
		categories = append(categories, c.ToModel())
	}

	nextPage := ""
	if o.LastEvaluatedKey["GSISK"] != nil {
		nextPage = util.FromDynamoID(*o.LastEvaluatedKey["GSISK"].S)
	}

	return categories, nextPage, nil
}
