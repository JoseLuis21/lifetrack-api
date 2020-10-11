package querybuilder

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"
)

// CategoryDynamo constructs an AWS DynamoDB query using a middleware strategy/pattern
type CategoryDynamo struct {
	query []expression.ConditionBuilder
}

// NewCategoryDynamo constructs an AWS DynamoDB query using a middleware strategy/pattern
//
// This factory method also assigns the given schema to the associated Primary Key
func NewCategoryDynamo(schema string) *CategoryDynamo {
	return &CategoryDynamo{
		query: []expression.ConditionBuilder{expression.Name("PK").BeginsWith(schema)},
	}
}

func (c *CategoryDynamo) User(username string) *CategoryDynamo {
	if username != "" {
		c.query = append(c.query, expression.Equal(expression.Name("GSIPK"), expression.Value(
			util.GenerateDynamoID("User", username))))
	}

	return c
}

func (c *CategoryDynamo) Title(title string) *CategoryDynamo {
	if title != "" {
		c.query = append(c.query, expression.Equal(expression.Name("title"), expression.Value(c.Title)))
	}

	return c
}

func (c *CategoryDynamo) Query(keyword string) *CategoryDynamo {
	if keyword != "" {
		c.query = append(c.query, expression.Contains(expression.Name("title"), keyword))
	}

	return c
}

func (c CategoryDynamo) Build() (expression.Expression, error) {
	// The following algorithm is required by the aws dynamodb library to ensure
	// correct mathematical boolean expressions
	if cLength := len(c.query); cLength >= 1 {
		if cLength == 2 {
			return expression.NewBuilder().WithFilter(expression.And(c.query[0], c.query[1])).Build()
		} else if cLength >= 3 {
			return expression.NewBuilder().WithFilter(expression.And(c.query[0], c.query[1], c.query[2:]...)).Build()
		}

		return expression.NewBuilder().WithFilter(c.query[0]).Build()
	}

	return expression.NewBuilder().Build()
}
