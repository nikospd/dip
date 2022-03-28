package utils

import (
	"reflect"
	"time"
)

type Application struct {
	AppId           string    `json:"appId" bson:"_id,omitempty"`
	UserId          string    `json:"userId" bson:"user_id,omitempty"`
	Description     string    `json:"description" bson:"description,omitempty"`
	CreatedAt       time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt      time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
	PersistRaw      bool      `json:"persistRaw" bson:"persist_raw,omitempty"`
	RawStorageId    string    `json:"rawStorageId" bson:"raw_storage_id,omitempty"`
	HasIntegrations bool      `json:"hasIntegrations" bson:"has_integrations,omitempty"`
	HasAutomations  bool      `json:"hasAutomations" bson:"has_automations,omitempty"`
	HasDevices      bool      `json:"hasDevices" bson:"has_devices,omitempty"`
	/*
		Future purposes: {
		HasDevices, DevicesIdPath, DataModel, AggregationRecipes, ShareDataWith
		}
	*/
}

type dict map[string]interface{}

func (d dict) d(k string) dict {
	return d[k].(map[string]interface{})
}

func (d dict) s(k string) string {
	return d[k].(string)
}

type IncomingMessage struct {
	Payload   map[string]interface{} `json:"payload" bson:"payload"`
	UserId    string                 `json:"userId" bson:"user_id"`
	AppId     string                 `json:"appId" bson:"app_id"`
	ArrivedAt time.Time              `json:"arrivedAt" bson:"arrived_at"`
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
			filteredObj[key] = data[key]
		}
	}
	return filteredObj
}
