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

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	secretsProvider := awsimpl.NewAWSSecretsProvider(cfg)

	secrets, err := secretsProvider.GetSecrets(ctx, []string{
		"gha-instrumentor/webhook-secret",
	})

	if err != nil {
		log.Println(err)
	}

	lambda.Start(
		otellambda.InstrumentHandler(func(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
			return awsimpl.Handler(ctx, request, secrets)
		}, xrayconfig.WithRecommendedOptions(tp)...),
	)
}
