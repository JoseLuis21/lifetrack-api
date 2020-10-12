package category

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"
)

// BuilderUserDynamo constructs and executes an AWS DynamoDB query fetching user's categories only
//	This is a concrete CategoryDynamo strategy implementation
type BuilderUserDynamo struct {
	user   string
	schema string
	exp    expression.Builder
	input  *dynamodb.QueryInput
}

func NewBuilderUserDynamo(tableName, indexName, schema string) *BuilderUserDynamo {
	return &BuilderUserDynamo{
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

func (b *BuilderUserDynamo) GetInput() *dynamodb.QueryInput {
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

func (b *BuilderUserDynamo) ByUser(user string) *BuilderUserDynamo {
	if user != "" {
		b.user = user
		b.exp = b.exp.WithKeyCondition(expression.KeyAnd(expression.Key("GSIPK").
			Equal(expression.Value(util.GenerateDynamoID("User", user))),
			expression.KeyBeginsWith(expression.Key("GSISK"), b.schema)))
	}

	return b
}

func (b *BuilderUserDynamo) Query(keyword string) *BuilderUserDynamo {
	if keyword != "" {
		b.exp = b.exp.WithFilter(expression.Contains(expression.Name("title"), keyword))
	}

	return b
}

func (b *BuilderUserDynamo) Limit(l int64) *BuilderUserDynamo {
	if l > 0 {
		b.input.SetLimit(l)
	}

	return b
}

// OrderBy true -> asc, false -> desc
func (b *BuilderUserDynamo) OrderBy(o bool) *BuilderUserDynamo {
	b.input.SetScanIndexForward(o)
	return b
}

func (b *BuilderUserDynamo) NextPage(token string) *BuilderUserDynamo {
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

func (b BuilderUserDynamo) Do(ctx context.Context, db *dynamodb.DynamoDB) ([]*model.Category, string, error) {
	o, err := db.QueryWithContext(ctx, b.GetInput())
	if err != nil {
		return nil, "", err
	}

	categories := make([]*model.Category, 0)
	for _, i := range o.Items {
		c := new(DynamoModel)
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
