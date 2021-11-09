package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	CreatedAt    time.Time `json:"createdAt" bson:"created_at,omitempty"`
}

func (t *PullSourceTask) HandleTask(col *mongo.Collection) {
	timeErr := time.Now().Sub(t.NextExecution).Milliseconds()
	t.LastExecuted = time.Now()
	t.NextExecution = time.Now().Add(time.Duration(t.Interval) * time.Second)
	t.Description = "Updated!!"
	updateQuery := bson.D{{"$set", t}}
	col.UpdateOne(context.TODO(), bson.D{{"_id", t.TaskId}}, updateQuery)
	fmt.Println("Handle task: ", t.TaskId, " with time error of: ", timeErr)
}

//TODO: Handle source info and actual execute the task (pull data) and add them into the rabbitmq

func (t PullSourceTask) ExecuteTask() {
	time.Sleep(3 * time.Second)
	fmt.Println("Execute for task: ", t.TaskId)
}
