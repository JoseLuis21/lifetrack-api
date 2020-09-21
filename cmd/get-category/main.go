package main

import (
	"context"
	"encoding/json"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/alexandria-oss/common-go/httputil"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"log"
	"net/http"
)

var muxLambda *gorillamux.GorillaMuxAdapter

func init() {
	q, err := dep.InjectGetCategoryQuery()
	if err != nil {
		log.Fatalf("failed to start get category query: %s", exception.GetDescription(err))
		return
	}

	log.Print("mux cold start")
	r := mux.NewRouter()
	r.Path("/live/category/{id}").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		category, err := q.Query(context.Background(), mux.Vars(r)["id"])
		if err != nil {
			httputil.RespondErrorJSON(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(category)
	})

	muxLambda = gorillamux.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%+v", req)
	return muxLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
