package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpPostIntegration struct {
	Uri string `json:"uri" bson:"uri,omitempty"`
}

func (i HttpPostIntegration) Send(msg IncomingMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(i.Uri, "application/json", bytes.NewBuffer(body))
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
