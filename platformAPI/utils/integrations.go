package utils

import (
	"errors"
)

const (
	HttpPost IntegrationTypes = "httpPost"
)

func (i *Integration) CheckType() error {
	switch i.IntegrationType {
	case HttpPost:
		return nil
	default:
		return errors.New("unsupported integration type")
	}
}

func (i HttpPostIntegration) Send(msg string) error {
	return nil
}

func (i HttpPostIntegration) CheckOption() error {
	if i.Uri == "" {
		return errors.New("http integration must include at least a uri")
	}
	return nil
}
