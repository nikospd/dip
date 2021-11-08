package utils

import (
	"fmt"
	"time"
)

type SourceTokenClaims struct {
	UserId      string    `json:"userId" bson:"user_id,omitempty"`
	AppId       string    `json:"appId" bson:"app_id,omitempty"`
	Description string    `json:"description" bson:"description,omitempty"`
	SourceToken string    `json:"sourceToken" bson:"_id,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt  time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
}

type LoginUserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   string `json:"userId" bson:"_id"`
}

func (s *SourceTokenClaims) PrintTheClaims() {
	fmt.Println("UserId: ", s.UserId)
	fmt.Println("AppId: ", s.AppId)
	fmt.Println("SourceToken: ", s.SourceToken)
	fmt.Println("CreatedAt: ", s.CreatedAt)
}

type Application struct {
	AppId         string    `json:"appId" bson:"_id,omitempty"`
	UserId        string    `json:"userId" bson:"user_id,omitempty"`
	Description   string    `json:"description" bson:"description,omitempty"`
	SourceType    string    `json:"sourceType" bson:"source_type,omitempty"`
	SourceDetails string    `json:"sourceDetails" bson:"source_details,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt    time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
	PersistRaw    bool      `json:"persistRaw" bson:"persist_raw,omitempty"`
	RawStorageId  string    `json:"rawStorageId" bson:"raw_storage_id,omitempty"`
	HasDevices    bool      `json:"hasDevices" bson:"has_devices,omitempty"`
	/*
		Future purposes: {
		HasDevices, DevicesIdPath, DataModel, AggregationRecipes, ShareDataWith
		}
	*/
}

type Storage struct {
	StorageId    string    `json:"storageId" bson:"_id,omitempty"`
	UserId       string    `json:"userId" bson:"user_id,omitempty"`
	AppId        string    `json:"appId" bson:"app_id,omitempty"`
	Type         string    `json:"type" bson:"type,omitempty"` //Cloud MongoDB, Proprietary DB, etc.
	Shared       bool      `json:"shared" bson:"shared,omitempty"`
	SharedWithId []string  `json:"sharedWithId" bson:"shared_with_id,omitempty"` //other user id
	Description  string    `json:"description" bson:"description,omitempty"`
	CreatedAt    time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt   time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
	/*
		Future purposes: {
			StorageConnectOptions string
			Wrong type. Should be something generic that depends on the type of the storage
		}
	*/
}

type UserResourcesStatus struct {
	UserId               string   `json:"userId" bson:"_id,omitempty"`
	SharedStoragesWithMe []string `json:"sharedStorageWithMe" bson:"shared_storage_with_me,omitempty"` //StorageIds that other users sharing with me
}
