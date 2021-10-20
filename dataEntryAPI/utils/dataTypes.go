package utils

import "time"

type SourceTokenClaims struct {
	UserId      string    `json:"userId" bson:"user_id"`
	AppId       string    `json:"appId" bson:"app_id"`
	Description string    `json:"description" bson:"description"`
	SourceToken string    `json:"sourceToken" bson:"_id"`
	CreatedAt   time.Time `json:"createdAt" bson:"created_at"`
	ModifiedAt  time.Time `json:"modifiedAt" bson:"modified_at"`
}

type IncomingMessage struct {
	Payload   map[string]interface{} `json:"payload"`
	UserId    string                 `json:"userId"`
	AppId     string                 `json:"appId"`
	ArrivedAt time.Time              `json:"arrivedAt"`
}
