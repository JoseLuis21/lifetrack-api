package main

import (
	"context"
	"log"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/lifetrack-api/pkg/dep"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/categoryhandler"
)

var muxLambda *gorillamux.GorillaMuxAdapter
var cleaning func()

func init() {
	cmd, clean, err := dep.InjectEditCategory()
	if err != nil {
		log.Fatalf("failed to start edit category command: %s", exception.GetDescription(err))
	}
	cleaning = clean

	h := categoryhandler.NewEdit(cmd, mux.NewRouter().PathPrefix("/live").Subrouter())
	log.Print("category handler successfully started")
	muxLambda = gorillamux.New(h.GetRouter())
}

// Handler AWS Lambda category handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.ProxyWithContext(ctx, req)
}

func main() {
	defer cleaning()
	lambda.Start(Handler)
}
