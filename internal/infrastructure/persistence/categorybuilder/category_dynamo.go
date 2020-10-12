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

// CategoryDynamo constructs and executes an AWS DynamoDB query fetching all categories
//	This is the default concrete CategoryDynamo strategy implementation
type CategoryDynamo struct {
	exp    expression.Builder
	schema string
	input  *dynamodb.ScanInput
}

func NewCategoryDynamo(table, schema string) *CategoryDynamo {
	return &CategoryDynamo{
		exp:    expression.NewBuilder(),
		schema: schema,
		input: &dynamodb.ScanInput{
			Limit:     aws.Int64(100),
			TableName: aws.String(table),
		},
	}
}

func (b *CategoryDynamo) GetInput() *dynamodb.ScanInput {
	b.exp = b.exp.WithFilter(expression.BeginsWith(expression.Name("PK"), b.schema))
	exp, _ := b.exp.Build()
	b.input.SetExpressionAttributeNames(exp.Names())
	b.input.SetExpressionAttributeValues(exp.Values())
	if exp.Filter() != nil {
		b.input.SetFilterExpression(*exp.Filter())
	}

	return b.input
}

func (b *CategoryDynamo) Query(keyword string) *CategoryDynamo {
	if keyword != "" {
		b.exp = b.exp.WithFilter(expression.Contains(expression.Name("title"), keyword))
	}

	return b
}

func (b *CategoryDynamo) Limit(l int64) *CategoryDynamo {
	if l > 0 && l <= 100 {
		b.input.SetLimit(l)
	}

	return b
}

func (b *CategoryDynamo) NextPage(token string) *CategoryDynamo {
	if token != "" {
		b.input.SetExclusiveStartKey(map[string]*dynamodb.AttributeValue{
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

func (b CategoryDynamo) Do(ctx context.Context, db *dynamodb.DynamoDB) ([]*model.Category, string, error) {
	o, err := db.ScanWithContext(ctx, b.GetInput())
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
	if o.LastEvaluatedKey["PK"] != nil {
		nextPage = util.FromDynamoID(*o.LastEvaluatedKey["PK"].S)
	}

	return categories, nextPage, nil
}
