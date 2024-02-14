package discord

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aJuvan/release-notify/notify/channels"
)

type DiscordWebhookPayload struct {
	Content string `json:"content"`
}

type DiscordChannel struct {
	Webhook string
}

func NewChannel(webhook string) channels.Channel {
	return &DiscordChannel{Webhook: webhook}
}

func (d *DiscordChannel) SendMessages(messages []string) error {
	for _, message := range messages {
		payload := DiscordWebhookPayload{Content: message}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", d.Webhook, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
	}

	return nil
}
