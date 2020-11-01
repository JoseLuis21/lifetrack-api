package serverless

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// StartLambda starts an AWS Lambda function with the given router
func StartLambda(lc fx.Lifecycle, logger *zap.Logger, h Handler) {
	muxLambda := gorillamux.New(h.GetRouter())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting aws lambda function")
			lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
				return muxLambda.ProxyWithContext(ctx, req)
			})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping aws lambda function")
			return nil
		},
	})
}
