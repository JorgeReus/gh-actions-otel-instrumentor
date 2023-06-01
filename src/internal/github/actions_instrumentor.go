package github

import (
	"context"
	"errors"
	"fmt"
	otelimpl "instrumentor/app/otel"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

)

func verifySignature(payloadBody string, secretToken string, signature string) error {
	hash := hmac.New(sha256.New, []byte(secretToken))
	hash.Write([]byte(payloadBody))
	expectedSignature := "sha256=" + hex.EncodeToString(hash.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
		return errors.New("Request signatures didn’t match!")
	}
	return nil
}

func Instrument(lambdaCtx context.Context, body string, signatureHeader string, webhookSecret string) error {
	var webhook WorkflowJobWebhook
	err := json.Unmarshal([]byte(body), &webhook)

	if err != nil {
		return errors.New("Error parsing webhook body")
	}
	if webhook.WorkflowJob.Status != statusCompleted.String() {
		return errors.New("Only completed events are processed")
	}

	err = verifySignature(body, webhookSecret, signatureHeader)

	if err != nil {
		return errors.New("Signatures don't match")
	}

	parentCtx := context.Background()

	// Instrument the lambda itself
	_, instrumentWorkflowSpan := otelimpl.GetTracerInstance().Start(lambdaCtx,
		fmt.Sprintf("Instrument workflow %s", webhook.WorkflowJob.Name),
	)
	defer instrumentWorkflowSpan.End()

	// Instrument the workflow in another trace
	workflowJobCtx, workflowJobSpan := otelimpl.GetTracerInstance().Start(parentCtx, webhook.WorkflowJob.Name, trace.WithTimestamp(webhook.WorkflowJob.StartedAt))
	workflowJobSpan.SetAttributes(
		attribute.Int64("github.resource.run_id", webhook.WorkflowJob.RunID),
		attribute.String("github.resource.html_url", webhook.WorkflowJob.HTMLURL),
		attribute.String("github.resource.runner_name", webhook.WorkflowJob.RunnerName),
		attribute.String("github.resource.head_branch", webhook.WorkflowJob.HeadBranch),
		attribute.StringSlice("github.resource.labels", webhook.WorkflowJob.Labels),
	)

	if webhook.WorkflowJob.Conclusion != conclusionSuccess.String() {
		workflowJobSpan.SetStatus(codes.Error, fmt.Sprintf("Workflowjob conclusion was %s", webhook.WorkflowJob.Conclusion))
	}

	defer workflowJobSpan.End(trace.WithTimestamp(webhook.WorkflowJob.CompletedAt))

	for _, step := range webhook.WorkflowJob.Steps {
		_, span := otelimpl.GetTracerInstance().Start(workflowJobCtx, step.Name, trace.WithTimestamp(step.StartedAt))
		span.SetAttributes(attribute.Int("number", step.Number))

		if step.Conclusion == conclusionFailure.String() || step.Conclusion == conclusionCancelled.String() {
			span.SetStatus(codes.Error, fmt.Sprintf("Step conclusion was %s", step.Conclusion))
		} else if step.Conclusion == conclusionSuccess.String() {
			span.SetStatus(codes.Ok, "")
		}

		span.End(trace.WithTimestamp(step.CompletedAt))
	}

	return nil
}
