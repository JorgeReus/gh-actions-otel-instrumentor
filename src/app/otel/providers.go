package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ProviderCancelFunc = func(context.Context) error

func InitTracerProvider(ctx context.Context) (ProviderCancelFunc, *trace.TracerProvider, error) {

	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating tracer provider: %w", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("Error shutting down tracer provider: %w", err)
		}
		return nil
	}, tp, nil
}
