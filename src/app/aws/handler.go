package aws

import (
	"context"
	"instrumentor/internal/github"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func handlerRespose(msg string, statusCode int) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		Body:       msg,
		StatusCode: statusCode,
	}
}

func Handler(lambdaCtx context.Context, request events.LambdaFunctionURLRequest, secrets map[string]string) (events.LambdaFunctionURLResponse, error) {
	signatureHeader := request.Headers["x-hub-signature-256"]

	if signatureHeader == "" {
		return handlerRespose("x-hub-signature-256 header not found", http.StatusBadRequest), nil
	}

	err := github.Instrument(lambdaCtx, request.Body, signatureHeader, secrets["gha-instrumentor/webhook-secret"])

	if err != nil {
		return handlerRespose(err.Error(), http.StatusBadRequest), nil
	}

	return handlerRespose("Workflow successfully instrumented", 200), nil
}
