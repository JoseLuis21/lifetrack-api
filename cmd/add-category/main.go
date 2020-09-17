package main

import (
	"context"
	"encoding/json"
	"github.com/alexandria-oss/common-go/httputil"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
	"log"
	"net/http"
)

var muxLambda *gorillamux.GorillaMuxAdapter

func init() {
	log.Print("mux cold start")
	r := mux.NewRouter()
	r.Path("/category").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, err := infrastructure.NewConfiguration()
		if err != nil {
			panic(err)
		}

		cmd := command.NewAddCategoryHandler(persistence.NewCategoryDynamoRepository(cfg))
		err = cmd.Handle(command.AddCategory{
			Ctx:         r.Context(),
			Title:       r.PostForm.Get("title"),
			User:        r.PostForm.Get("user"),
			Description: r.PostForm.Get("description"),
		})

		if err != nil {
			httputil.RespondErrorJSON(err, w)
			return
		}

		_ = json.NewEncoder(w).Encode(httputil.Response{
			Message: "successfully created category",
			Code:    http.StatusOK,
		})
	})

	muxLambda = gorillamux.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
