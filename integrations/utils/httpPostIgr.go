package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type HttpPostIntegration struct {
	Uri     string            `json:"uri" bson:"uri,omitempty"`
	Headers map[string]string `json:"headers" bson:"headers,omitempty"`
}

func (i HttpPostIntegration) Send(msg IncomingMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	timeout := time.Duration(3 * time.Second)
	client := http.Client{Timeout: timeout}

	request, err := http.NewRequest("POST", i.Uri, bytes.NewBuffer(body))
	request.Header.Set("content-type", "application/json")
	for k, v := range i.Headers {
		request.Header.Set(k, v)
	}
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func (i HttpPostIntegration) CheckOption() error {
	return nil
}
