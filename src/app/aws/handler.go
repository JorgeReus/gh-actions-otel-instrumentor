package aws

import (
	"context"
	"fmt"
	"instrumentor/app/discord"
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

	notificationUrl := secrets["gha-instrumentor/notification-url"]

	if signatureHeader == "" {
		return handlerRespose("x-hub-signature-256 header not found", http.StatusBadRequest), nil
	}

	ghWebhook, err := github.Instrument(lambdaCtx, request.Body, signatureHeader, secrets["gha-instrumentor/webhook-secret"])

	if err != nil {
		return handlerRespose(err.Error(), http.StatusBadRequest), nil
	}

	conclusionColor := discord.DiscordColorGreen

	if ghWebhook.WorkflowJob.Conclusion != "success" {
		conclusionColor = discord.DiscordColorRed
	}
	dsWebhook := discord.DiscordWebhook{
		Username:  "Github Workflow Notifier",
		AvatarURL: "https://avatars.githubusercontent.com/u/45807407?v=4",
		Content:   fmt.Sprintf("Workflow %s from repo %s has finished", ghWebhook.WorkflowJob.WorkflowName, ghWebhook.Repository.Name),
		Embeds: []discord.EmbedObject{
			{
				Author: discord.AuthorObject{
					Name:    ghWebhook.Sender.Login,
					URL:     ghWebhook.Sender.HTMLURL,
					IconURL: ghWebhook.Sender.AvatarURL,
				},
				Title:       ghWebhook.WorkflowJob.Name,
				URL:         ghWebhook.WorkflowJob.HTMLURL,
				Description: "Workflow job webhook details",
				Color:       discord.DiscordColorBlue,
				Fields: []discord.FieldObject{
					{
						Name:   "Started At",
						Value:  ghWebhook.WorkflowJob.StartedAt.String(),
						Inline: true,
					},
					{
						Name:   "Completed At",
						Value:  ghWebhook.WorkflowJob.CompletedAt.String(),
						Inline: true,
					},
					{
						Name:   "Runner Name",
						Value:  ghWebhook.WorkflowJob.RunnerName,
						Inline: false,
					},
				},
			},
			{
				Title: fmt.Sprintf("Conclusion: %s", ghWebhook.WorkflowJob.Conclusion),
				Color: conclusionColor,
			},
		},
	}
	err = discord.Notify(notificationUrl, dsWebhook)
	if err != nil {
		return handlerRespose(err.Error(), http.StatusBadRequest), nil
	}

	return handlerRespose("Workflow successfully instrumented", 200), nil
}
