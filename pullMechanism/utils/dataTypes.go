package utils

import (
	"time"
)

type PullSourceTask struct {
	TaskId        string    `json:"taskId" bson:"_id,omitempty"`
	UserId        string    `json:"userId" bson:"user_id,omitempty"`
	AppId         string    `json:"appId" bson:"app_id,omitempty"`
	SourceURI     string    `json:"sourceURI" bson:"source_uri,omitempty"`
	Interval      int       `json:"interval" bson:"interval,omitempty"`
	Description   string    `json:"description" bson:"description,omitempty"`
	LastExecuted  time.Time `json:"lastExecuted" bson:"last_executed,omitempty"`
	NextExecution time.Time `json:"nextExecution" bson:"next_execution,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"created_at,omitempty"`
}

type IncomingMessage struct {
	Payload   map[string]interface{} `json:"payload"`
	UserId    string                 `json:"userId"`
	AppId     string                 `json:"appId"`
	ArrivedAt time.Time              `json:"arrivedAt"`
}
