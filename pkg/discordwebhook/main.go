package discordwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// WebhookPayload represents the structure of the Discord webhook payload.
type WebhookPayload struct {
	Content   string  `json:"content,omitempty"`
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

// Embed represents a Discord embed object.
type Embed struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	URL         string  `json:"url,omitempty"`
	Color       int     `json:"color,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
	Footer      *Footer `json:"footer,omitempty"`
}

// Field represents a field within a Discord embed.
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Footer represents the footer in a Discord embed.
type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// SendWebhook sends a message or embed to a Discord webhook URL.
func SendWebhook(webhookURL string, payload WebhookPayload) error {
	// Marshal the payload into JSON.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create an HTTP POST request.
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set appropriate headers.
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP response errors.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("received non-successful HTTP status code: %d", resp.StatusCode))
	}

	return nil
}
