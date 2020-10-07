package main

import (
	"context"
	"log"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"github.com/neutrinocorp/life-track-api/pkg/transport/handler"
)

var muxLambda *gorillamux.GorillaMuxAdapter
var cleaning func()

func init() {
	cmd, clean, err := dep.InjectAddCategoryHandler()
	if err != nil {
		log.Fatalf("failed to start add category command: %s", exception.GetDescription(err))
	}
	cleaning = clean

	h := handler.NewAddCategory(cmd, mux.NewRouter())
	log.Print("handler successfully started")
	muxLambda = gorillamux.New(h.GetRouter())
}

// Handler AWS Lambda handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.ProxyWithContext(ctx, req)
}

func main() {
	defer cleaning()
	lambda.Start(Handler)
}
