package main

import (
	"context"
	"dev.com/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var client *mongo.Client
var channel *amqp.Channel
var queue amqp.Queue

func main() {
	/*
		Connect to MongoDB
	*/
	clientOptions := options.Client().ApplyURI("mongodb://test:test@localhost:27017/")
	var connectionError error
	client, connectionError = mongo.Connect(context.TODO(), clientOptions)
	if connectionError != nil {
		log.Fatalln(connectionError)
	}
	/*
		Connect to RabbitMQ Server
	*/
	conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	queue, err = channel.QueueDeclare(
		"incoming_data", true, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	/*
		Prepare the web server
	*/
	e := echo.New()
	/*
		/*
			Assign resources to unauthenticated endpoints
	*/
	e.POST("/data", dataEntry)
	/*
		Start the web server
	*/
	e.Logger.Fatal(e.Start(":1324"))
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
	}
}

//TODO: Set a dead letter queue
//TODO: Set up configuration files and a proper logger

func dataEntry(c echo.Context) error {
	sourceToken := c.Request().Header.Get("source-token")
	if sourceToken == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"msg": "Source token not provided"})
	}
	collection := client.Database("staging").Collection("source_tokens")
	cur := collection.FindOne(context.TODO(), bson.D{
		{"_id", sourceToken}})
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusUnauthorized, echo.Map{"msg": "The source token that provided is unauthorized"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error bad gateway"})
	}
	var stc utils.SourceTokenClaims
	err := cur.Decode(&stc)
	if err != nil {
		failOnError(err, "Failed to decode token claims")
	}
	var body map[string]interface{}
	err = c.Bind(&body)
	if err != nil {
		failOnError(err, "Failed to bind request body")
	}
	var msg utils.IncomingMessage
	msg.Payload = body
	msg.UserId = stc.UserId
	msg.AppId = stc.AppId
	msg.ArrivedAt = time.Now()
	msgJs, _ := json.Marshal(msg)
	err = channel.Publish("", queue.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msgJs,
		})
	if err != nil {
		failOnError(err, "Failed to Publish message")
	}
	return c.JSON(http.StatusAccepted, echo.Map{"msg": "OK"})
}
