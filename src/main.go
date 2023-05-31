package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	otelimpl "instrumentor/app/otel"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func handlerRespose(msg string, statusCode int) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		Body:       msg,
		StatusCode: statusCode,
	}
}

func handler(lambdaCtx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	var webhook WorkflowJobWebhook
	err := json.Unmarshal([]byte(request.Body), &webhook)

	if err != nil {
		return handlerRespose("Error parsing webhook body", 400), nil
	}

	if webhook.WorkflowJob.Status != StatusCompleted.String() {
		return handlerRespose("Only completed events are processed", 400), nil
	}

	parentCtx := context.Background()

	// Instrument the lambda itself
	_, instrumentWorkflowSpan := otelimpl.GetTracerInstance().Start(lambdaCtx,
		fmt.Sprintf("Instrument workflow %s", webhook.WorkflowJob.Name),
	)
	defer instrumentWorkflowSpan.End()

	// Instrument the workflow in another trace
	ctx, jobSpan := otelimpl.GetTracerInstance().Start(parentCtx, webhook.WorkflowJob.Name, trace.WithTimestamp(webhook.WorkflowJob.StartedAt))
	jobSpan.SetAttributes(
		attribute.Int64("github.resource.run_id", webhook.WorkflowJob.RunID),
		attribute.String("github.resource.html_url", webhook.WorkflowJob.HTMLURL),
		attribute.String("github.resource.runner_name", webhook.WorkflowJob.RunnerName),
		attribute.String("github.resource.head_branch", webhook.WorkflowJob.HeadBranch),
		attribute.StringSlice("github.resource.labels", webhook.WorkflowJob.Labels),
	)

	if webhook.WorkflowJob.Conclusion != ConclusionSuccess.String() {
		jobSpan.SetStatus(codes.Error, fmt.Sprintf("Workflowjob conclusion was %s", webhook.WorkflowJob.Conclusion))
	}

	defer jobSpan.End(trace.WithTimestamp(webhook.WorkflowJob.CompletedAt))

	for _, step := range webhook.WorkflowJob.Steps {
		_, span := otelimpl.GetTracerInstance().Start(ctx, step.Name, trace.WithTimestamp(step.StartedAt))
		span.SetAttributes(attribute.Int("number", step.Number))

		if step.Conclusion == ConclusionFailure.String() || step.Conclusion == ConclusionCancelled.String() {
			span.SetStatus(codes.Error, fmt.Sprintf("Step conclusion was %s", step.Conclusion))
		} else if step.Conclusion == ConclusionSuccess.String() {
			span.SetStatus(codes.Ok, "")
		}

		span.End(trace.WithTimestamp(step.CompletedAt))
	}

	return handlerRespose("Respons3", 200), nil
}

func main() {
	ctx := context.Background()
	providerCancelFun, tp, err := otelimpl.InitTracerProvider(ctx)

	if err != nil {
		log.Fatalln(fmt.Errorf("Cannot start otel tracer: %w", err))
	}

	defer providerCancelFun(ctx)

	lambda.Start(
		otellambda.InstrumentHandler(handler, xrayconfig.WithRecommendedOptions(tp)...),
	)
}
