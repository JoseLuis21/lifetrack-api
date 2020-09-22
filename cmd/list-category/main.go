package main

import (
	"context"
	"encoding/json"
	"github.com/alexandria-oss/common-go/httputil"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"log"
	"net/http"
)

var muxAdapter *gorillamux.GorillaMuxAdapter

func init() {
	q, err := dep.InjectListCategoriesQuery()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.StrictSlash(true).Path("/live/category").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categories, token, err := q.Query(r.Context(), r.URL.Query().Get("next_token"), r.URL.Query().Get("page_size"), map[string]string{
			"user":  r.URL.Query().Get("u"),
			"title": r.URL.Query().Get("title"),
			"query": r.URL.Query().Get("q"),
		})
		if err != nil {
			httputil.RespondErrorJSON(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(struct {
			Categories []*model.Category `json:"categories"`
			TotalItems int               `json:"total_items"`
			NextToken  string            `json:"next_token"`
		}{
			Categories: categories,
			TotalItems: len(categories),
			NextToken:  token,
		})
	})

	muxAdapter = gorillamux.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxAdapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
