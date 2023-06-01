package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordWebhook struct {
	Username  string        `json:"username"`
	AvatarURL string        `json:"avatar_url"`
	Content   string        `json:"content"`
	Embeds    []EmbedObject `json:"embeds"`
}

type EmbedObject struct {
	Author      AuthorObject    `json:"author"`
	Title       string          `json:"title"`
	URL         string          `json:"url"`
	Description string          `json:"description"`
	Color       int             `json:"color"`
	Fields      []FieldObject   `json:"fields"`
	Thumbnail   ThumbnailObject `json:"thumbnail"`
	Image       ImageObject     `json:"image"`
	Footer      FooterObject    `json:"footer"`
}

type AuthorObject struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

type FieldObject struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type ThumbnailObject struct {
	URL string `json:"url"`
}

type ImageObject struct {
	URL string `json:"url"`
}

type FooterObject struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

const (
  DiscordColorRed = 16711680
  DiscordColorGreen = 14177041
  DiscordColorBlue = 1127128
)


func Notify(url string, webhook DiscordWebhook) error {
	jsonPayload, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("Error marshaling JSON payload: %w", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("Error sending request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected message. Expected: %d, Got: %d", http.StatusNoContent, res.StatusCode)
	}
	return nil
}
