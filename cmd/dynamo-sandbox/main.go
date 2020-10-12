package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/neutrinocorp/life-track-api/internal/domain/model"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/readmodel"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/awsutil"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
)

func main() {
	r := mux.NewRouter()
	r.Path("/sandbox/user/{user}/category").Methods(http.MethodGet).HandlerFunc(fetchCategories)

	log.Print("Listening HTTP TCP socket at :8080")
	panic(http.ListenAndServe(":8080", r))
}

func fetchCategories(w http.ResponseWriter, r *http.Request) {
	db := persistence.NewDynamoConn(awsutil.NewSession(), "us-east-1")

	var err error
	limit := int64(100)
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			httputil.RespondErrorJSON(err, w)
			return
		}
	}

	order := true
	if o := r.URL.Query().Get("order_by"); strings.ToUpper(o) == "DESC" {
		order = false
	}

	categories, nextPage, err := NewInputBuilder().ByUser(mux.Vars(r)["user"]).
		Query(r.URL.Query().Get("query")).Limit(limit).OrderBy(order).NextPage(r.URL.Query().
		Get("next_page")).Do(db)
	if err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Categories []*model.Category `json:"categories"`
		NextPage   string            `json:"next_page"`
	}{
		Categories: categories,
		NextPage:   nextPage,
	})
}

type InputBuilder struct {
	user  string
	exp   expression.Builder
	input *dynamodb.QueryInput
}

func NewInputBuilder() *InputBuilder {
	return &InputBuilder{
		user: "",
		exp:  expression.NewBuilder(),
		input: &dynamodb.QueryInput{
			IndexName: aws.String("GSIPK-index"),
			Limit:     aws.Int64(100),
			TableName: aws.String("lifetrack-dev"),
		},
	}
}

func (b *InputBuilder) GetInput() *dynamodb.QueryInput {
	exp, _ := b.exp.Build()
	b.input.ExpressionAttributeValues = exp.Values()
	b.input.ExpressionAttributeNames = exp.Names()
	b.input.FilterExpression = exp.Filter()
	b.input.KeyConditionExpression = exp.KeyCondition()
	return b.input
}

func (b *InputBuilder) ByUser(user string) *InputBuilder {
	if user != "" {
		b.user = user
		b.exp = b.exp.WithKeyCondition(expression.KeyAnd(expression.Key("GSIPK").
			Equal(expression.Value(util.GenerateDynamoID("User", user))),
			expression.KeyBeginsWith(expression.Key("GSISK"), "Category")))
	}

	return b
}

func (b *InputBuilder) Query(keyword string) *InputBuilder {
	if keyword != "" {
		b.exp = b.exp.WithFilter(expression.Contains(expression.Name("title"), keyword))
	}

	return b
}

func (b *InputBuilder) Limit(l int64) *InputBuilder {
	if l > 0 {
		b.input.Limit = aws.Int64(l)
	}

	return b
}

// OrderBy true -> asc, false -> desc
func (b *InputBuilder) OrderBy(o bool) *InputBuilder {
	b.input.ScanIndexForward = aws.Bool(o)
	return b
}

func (b *InputBuilder) NextPage(token string) *InputBuilder {
	if token != "" {
		b.input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"GSIPK": {
				S: aws.String(util.GenerateDynamoID("User", b.user)),
			},
			"GSISK": {
				S: aws.String(util.GenerateDynamoID("Category", token)),
			},
			"PK": {
				S: aws.String(util.GenerateDynamoID("Category", token)),
			},
			"SK": {
				S: aws.String(util.GenerateDynamoID("Category", token)),
			},
		}
	}

	return b
}

func (b InputBuilder) Do(db *dynamodb.DynamoDB) ([]*model.Category, string, error) {
	o, err := db.Query(b.GetInput())
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
