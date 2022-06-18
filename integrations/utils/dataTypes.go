package utils

import (
	"errors"
	"reflect"
	"time"
)

type IncomingMessage struct {
	Payload   map[string]interface{} `json:"payload"`
	UserId    string                 `json:"userId,omitempty"`
	AppId     string                 `json:"appId,omitempty"`
	ArrivedAt time.Time              `json:"arrivedAt"`
}

type Integration struct {
	Id              string           `json:"id" bson:"_id,omitempty"`
	AppId           string           `json:"appId" bson:"app_id,omitempty"`
	UserId          string           `json:"userId" bson:"user_id,omitempty"`
	Description     string           `json:"description" bson:"description,omitempty"`
	IntegrationType IntegrationTypes `json:"type" bson:"type,omitempty"`
	//Change this to take IntegrationOption interface
	Option     HttpPostIntegration `json:"option" bson:"option,omitempty"`
	CreatedAt  time.Time           `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt time.Time           `json:"modifiedAt" bson:"modified_at,omitempty"`
}

type IntegrationTypes string

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

func (i *Integration) Send(msg IncomingMessage) error {
	return i.Option.Send(msg)
}

type IntegrationOption interface {
	Send(message IncomingMessage) error
	CheckOption() error
}

type StorageFilter struct {
	FilterId    string                 `json:"filterId" bson:"_id,omitempty"`
	UserId      string                 `json:"userId" bson:"user_id,omitempty"`
	StorageId   string                 `json:"storageId" bson:"storage_id,omitempty"`
	Description string                 `json:"description" bson:"description,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"`
	CreatedAt   time.Time              `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt  time.Time              `json:"modifiedAt" bson:"modified_at,omitempty"`
}

func (f *StorageFilter) Apply(msg *IncomingMessage) error {
	msg.Payload = filterHelper(msg.Payload, f.Attributes)
	return nil
}
func filterHelper(data map[string]interface{}, filter map[string]interface{}) map[string]interface{} {
	filteredObj := make(map[string]interface{})
	for key := range filter {
		if reflect.TypeOf(filter[key]) == reflect.TypeOf(filteredObj) {
			if _, ok := data[key]; ok {
				filteredObj[key] = filterHelper(data[key].(map[string]interface{}), filter[key].(map[string]interface{}))
			}
		} else {
			if filter[key].(float64) == 1 {
				filteredObj[key] = data[key]
			}
		}
	}
	return filteredObj
}
