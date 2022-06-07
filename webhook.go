package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Author struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      Author `json:"author"`
	Color       int    `json:"color"`
}

func queueRequest(a []Embed) {
	var b []Embed
	for i, p := range a {
		b = append(b, p)
		if (i+1)%10 == 0 {
			err := postWebhook(b)
			if err != nil {
				return
			}
			b = nil
		}
	}
	if len(b) > 0 {
		err := postWebhook(b)
		if err != nil {
			return
		}
	}
}

func postWebhook(e []Embed) error {
	b := Body{Embeds: e}
	jsonData, _ := json.Marshal(b)
	_, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}
