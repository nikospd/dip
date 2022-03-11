package utils

import (
	"time"
)

type Application struct {
	AppId        string    `json:"appId" bson:"_id,omitempty"`
	UserId       string    `json:"userId" bson:"user_id,omitempty"`
	Description  string    `json:"description" bson:"description,omitempty"`
	CreatedAt    time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt   time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
	PersistRaw   bool      `json:"persistRaw" bson:"persist_raw,omitempty"`
	RawStorageId string    `json:"rawStorageId" bson:"raw_storage_id,omitempty"`
	HasDevices   bool      `json:"hasDevices" bson:"has_devices,omitempty"`
	/*
		Future purposes: {
		HasDevices, DevicesIdPath, DataModel, AggregationRecipes, ShareDataWith
		}
	*/
}

type IncomingMessage struct {
	Payload   map[string]interface{} `json:"payload" bson:"payload"`
	UserId    string                 `json:"userId" bson:"user_id"`
	AppId     string                 `json:"appId" bson:"app_id"`
	ArrivedAt time.Time              `json:"arrivedAt" bson:"arrived_at"`
}

type StorageFilter struct {
	FilterId    string    `json:"filterId" bson:"_id,omitempty"`
	UserId      string    `json:"userId" bson:"user_id,omitempty"`
	StorageId   string    `json:"storageId" bson:"storage_id,omitempty"`
	Description string    `json:"description" bson:"description,omitempty"`
	Attributes  []string  `json:"attributes,omitempty" bson:"attributes,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt  time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
}

func (f *StorageFilter) Apply(msg *IncomingMessage) error {
	filtered_payload := make(map[string]interface{})
	for _, s := range f.Attributes {
		filtered_payload[s] = msg.Payload[s]
	}
	msg.Payload = filtered_payload
	return nil
}
