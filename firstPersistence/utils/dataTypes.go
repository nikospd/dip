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
	FilterId    string     `json:"filterId" bson:"_id,omitempty"`
	UserId      string     `json:"userId" bson:"user_id,omitempty"`
	StorageId   string     `json:"storageId" bson:"storage_id,omitempty"`
	Description string     `json:"description" bson:"description,omitempty"`
	Attributes  [][]string `json:"attributes,omitempty" bson:"attributes,omitempty"`
	CreatedAt   time.Time  `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt  time.Time  `json:"modifiedAt" bson:"modified_at,omitempty"`
}

func (f *StorageFilter) Apply(msg *IncomingMessage) error {
	filteredObj := make(map[string]interface{})
	for _, f1 := range f.Attributes {
		tmp := msg.Payload
		for _, f2 := range f1[:len(f1)-1] {
			if _, ok := tmp[f2]; ok {
				tmp = tmp[f2].(map[string]interface{})
			} else {
				break
			}
		}
		item := tmp[f1[len(f1)-1]]
		if item == nil {
			continue
		}
		for i := len(f1) - 1; i >= 1; i-- {
			tmp := map[string]interface{}{f1[i]: item}
			item = tmp
		}
		filteredObj[f1[0]] = item
	}
	msg.Payload = filteredObj
	return nil
}
