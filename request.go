package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func sendRequest(e []Embed) error {
	b := Body{Embeds: e}
	jsonData, _ := json.Marshal(b)
	_, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}
