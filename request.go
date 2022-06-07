package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func QueueRequest(a []Embed) {
	var b []Embed
	for i, p := range a {
		b = append(b, p)
		if (i+1)%10 == 0 {
			err := sendRequest(b)
			if err != nil {
				return
			}
			b = nil
		}
	}
	if len(b) > 0 {
		err := sendRequest(b)
		if err != nil {
			return
		}
	}
}

func sendRequest(e []Embed) error {
	b := Body{Embeds: e}
	jsonData, _ := json.Marshal(b)
	_, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}
