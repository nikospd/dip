package main

//TODO: Set a dead letter queue
//TODO: Set up configuration files and a proper logger

import (
	"context"
	"dev.com/utils"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

	msgs, err := channel.Consume(queue.Name, "", false, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Println("new incoming message")
			var msg utils.IncomingMessage
			err := json.Unmarshal(d.Body, &msg)
			failOnError(err, "Failed to read incoming message")
			//Search for the application
			var application utils.Application
			appCollection := client.Database("staging").Collection("applications")
			one := appCollection.FindOne(context.TODO(), bson.D{
				{"_id", msg.AppId}, {"user_id", msg.UserId}})
			if one.Err() != nil {
				if one.Err() == mongo.ErrNoDocuments {
					failOnError(one.Err(), "No application for those app and user id")
				} else {
					failOnError(one.Err(), "Failed on searching application from database")
				}
				//TODO: put the message in a dead letter queue
				err = d.Ack(false)
				failOnError(err, "Failed to acknowledge")
			}
			err = one.Decode(&application)
			failOnError(err, "Failed to decode application from database")
			//Save the message if raw persistence is activated
			if application.PersistRaw {
				fmt.Println("With raw persistence enabled")
				persistCollection := client.Database("staging_payloads").Collection(application.RawStorageId)
				persistCollection.InsertOne(context.TODO(), msg)
			}
			//Make the acknowledgment
			err = d.Ack(false)
			failOnError(err, "Failed to acknowledge")
		}
	}()
	fmt.Println("Start consuming...")
	<-forever
	fmt.Println("End of program")
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
	}
}
