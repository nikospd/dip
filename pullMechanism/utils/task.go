package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"time"
)

func (t PullSourceTask) ExecuteTask(channel *amqp.Channel, queue amqp.Queue) {
	resp, err := http.Get(t.SourceURI)
	if err != nil {
		fmt.Println(t.TaskId, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(t.TaskId, err)
		return
	}
	fmt.Println(t.TaskId, string(body))
	var msg IncomingMessage
	json.Unmarshal(body, &msg.Payload)
	msg.UserId = t.UserId
	msg.AppId = t.AppId
	msg.ArrivedAt = time.Now()
	fmt.Println(msg)
	msgJs, _ := json.Marshal(msg)
	err = channel.Publish("", queue.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msgJs,
		})
	FailOnError(err, "Failed to Publish message")
}

func (t *PullSourceTask) HandleTask(col *mongo.Collection) {
	timeErr := time.Now().Sub(t.NextExecution).Milliseconds()
	t.LastExecuted = time.Now()
	t.NextExecution = time.Now().Add(time.Duration(t.Interval) * time.Minute)
	t.Description = "Updated!!"
	updateQuery := bson.D{{"$set", t}}
	col.UpdateOne(context.TODO(), bson.D{{"_id", t.TaskId}}, updateQuery)
	fmt.Println("Handle task: ", t.TaskId, " with time error of: ", timeErr)
}
