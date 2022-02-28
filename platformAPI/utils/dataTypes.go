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

type User struct {
	Username   string    `json:"username" bson:"username,omitempty"`
	Password   string    `json:"password,omitempty" bson:"password,omitempty"`
	Email      string    `json:"email" bson:"email,omitempty"`
	UserId     string    `json:"userId" bson:"_id"`
	CreatedAt  time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
	LastLogin  time.Time `json:"lastLogin" bson:"last_login,omitempty"`
}

func (s *SourceTokenClaims) PrintTheClaims() {
	fmt.Println("UserId: ", s.UserId)
	fmt.Println("AppId: ", s.AppId)
	fmt.Println("SourceToken: ", s.SourceToken)
	fmt.Println("CreatedAt: ", s.CreatedAt)
}

type Application struct {
	AppId              string    `json:"appId" bson:"_id,omitempty"`
	UserId             string    `json:"userId" bson:"user_id,omitempty"`
	ApplicationGroupId string    `json:"applicationGroupId" bson:"application_group_id,omitempty"`
	Description        string    `json:"description" bson:"description,omitempty"`
	SourceType         string    `json:"sourceType" bson:"source_type,omitempty"`
	CreatedAt          time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt         time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
	PersistRaw         bool      `json:"persistRaw" bson:"persist_raw,omitempty"`
	RawStorageId       string    `json:"rawStorageId" bson:"raw_storage_id,omitempty"`
	/*
		Future purposes: {
		HasDevices, DevicesIdPath, DataModel, AggregationRecipes
		AggregationStorages: {AggregationId, StorageId}
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

type PullSourceTask struct {
	TaskId        string    `json:"taskId" bson:"_id,omitempty"`
	UserId        string    `json:"userId" bson:"user_id,omitempty"`
	AppId         string    `json:"appId" bson:"app_id,omitempty"`
	SourceURI     string    `json:"sourceURI" bson:"source_uri,omitempty"`
	Interval      int       `json:"interval" bson:"interval,omitempty"`
	Enabled       bool      `json:"enabled" bson:"enabled,omitempty"`
	Description   string    `json:"description" bson:"description,omitempty"`
	LastExecuted  time.Time `json:"lastExecuted" bson:"last_executed,omitempty"`
	NextExecution time.Time `json:"nextExecution" bson:"next_execution,omitempty"`
	ModifiedAt    time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"created_at,omitempty"`
}

type ApplicationGroup struct {
	GroupId           string    `json:"groupId" bson:"_id,omitempty"`
	UserId            string    `json:"userId" bson:"user_id,omitempty"`
	Description       string    `json:"description" bson:"description,omitempty"`
	Applications      []string  `json:"applications,omitempty" bson:"applications,omitempty"`
	NumOfApplications int       `json:"numOfApplications" bson:"num_of_applications,omitempty"`
	CreatedAt         time.Time `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt        time.Time `json:"modifiedAt" bson:"modified_at,omitempty"`
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
