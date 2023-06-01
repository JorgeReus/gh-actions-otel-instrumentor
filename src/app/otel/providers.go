package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitTracerProvider(ctx context.Context) (*trace.TracerProvider, error) {

	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating tracer provider: %w", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

  return tp, nil
}
