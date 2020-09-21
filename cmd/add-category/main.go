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
	"github.com/neutrinocorp/life-track-api/internal/application/command"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"log"
	"net/http"
)

var muxLambda *gorillamux.GorillaMuxAdapter

func init() {
	cmd, err := dep.InjectAddCategoryHandler()
	if err != nil {
		log.Fatalf("failed to start add category command: %s", exception.GetDescription(err))
		return
	}

	log.Print("mux cold start")
	r := mux.NewRouter()
	r.Path("/live/category").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = cmd.Handle(command.AddCategory{
			Ctx:         r.Context(),
			Title:       r.PostFormValue("title"),
			User:        r.PostFormValue("user"),
			Description: r.PostFormValue("description"),
		})

		if err != nil {
			httputil.RespondErrorJSON(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
