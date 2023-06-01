package otel

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setEnvVars() {
	_ = os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "testFunction")
	_ = os.Setenv("AWS_REGION", "us-texas-1")
	_ = os.Setenv("AWS_LAMBDA_FUNCTION_VERSION", "$LATEST")
	_ = os.Setenv("AWS_LAMBDA_LOG_STREAM_NAME", "2023/01/01/[$LATEST]5d1edb9e525d486696cf01a3503487bc")
	_ = os.Setenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE", "128")
	_ = os.Setenv("_X_AMZN_TRACE_ID", "Root=1-5759e988-bd862e3fe1be46a994272793;Parent=53995c3f42cd8ad8;Sampled=1")
}

func TestErrorTracerProvider(t *testing.T) {
	_, err := InitTracerProvider(context.Background())

	assert.Error(t, err)
}

func TestNoErrorTracerProvider(t *testing.T) {
	setEnvVars()
	mockCollector := runMockCollectorAtEndpoint(t, ":4317")
	defer func() {
		_ = mockCollector.Stop()
	}()
	<-time.After(5 * time.Millisecond)
	tp, err := InitTracerProvider(context.Background())

  GetTracerInstance()


  defer func() error {
		if err := tp.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("Error shutting down tracer provider: %w", err)
		}
		return nil
	}()

	assert.NoError(t, err)
}
