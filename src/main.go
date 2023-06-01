package main

import (
	"context"
	"fmt"
	otelimpl "instrumentor/app/otel"

	awsimpl "instrumentor/app/aws"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func main() {
	ctx := context.Background()

	tp, err := otelimpl.InitTracerProvider(ctx)

	if err != nil {
		log.Fatalln(fmt.Errorf("Cannot start otel tracer: %w", err))
	}

	defer func() error {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("Error shutting down tracer provider: %w", err)
		}
		return nil
	}()

	lambda.Start(
		otellambda.InstrumentHandler(func(lambdaCtx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
			cfg, err := config.LoadDefaultConfig(lambdaCtx)
			if err != nil {
				log.Fatal(err)
			}
	    otelaws.AppendMiddlewares(&cfg.APIOptions)

			secretsProvider := awsimpl.NewAWSSecretsProvider(cfg)

			secrets, err := secretsProvider.GetSecrets(lambdaCtx, []string{
				"gha-instrumentor/webhook-secret",
				"gha-instrumentor/notification-url",
			})

			if err != nil {
				log.Fatal(err)
			}
			return awsimpl.Handler(lambdaCtx, request, secrets)
		}, xrayconfig.WithRecommendedOptions(tp)...),
	)
}
