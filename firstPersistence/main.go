package main

//TODO: Set a dead letter queue
//TODO: Set up configuration files and a proper logger

import (
	"context"
	"dev.com/config"
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
var cfg config.Configuration

func main() {
	/*
		Read configuration file
	*/
	config.ReadConf("config.json", &cfg)
	/*
		Connect to MongoDB
	*/
	mongoCred := cfg.MongoCredentials
	mongoUri := utils.MongoCredentials(mongoCred.User, mongoCred.Password, mongoCred.Host, mongoCred.Port)
	clientOptions := options.Client().ApplyURI(mongoUri)
	var connectionError error
	client, connectionError = mongo.Connect(context.TODO(), clientOptions)
	if connectionError != nil {
		log.Fatalln(connectionError)
	}
	/*
		Connect to RabbitMQ Server
	*/
	amqpCred := cfg.AmqpCredentials
	amqpUri := utils.AmqpCredentials(amqpCred.User, amqpCred.Password, amqpCred.Host, amqpCred.Port)
	conn, err := amqp.Dial(amqpUri)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	channel, err = conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	queue, err = channel.QueueDeclare(
		cfg.AmqpQueues.IncomingData, true, false, false, false, nil)
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(queue.Name, "", false, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg utils.IncomingMessage
			err := json.Unmarshal(d.Body, &msg)
			utils.FailOnError(err, "Failed to read incoming message")
			//Search for the application
			var application utils.Application
			appCollection := client.Database(cfg.MongoDatabase.Resources).Collection(cfg.MongoCollection.Applications)
			one := appCollection.FindOne(context.TODO(), bson.D{
				{"_id", msg.AppId}, {"user_id", msg.UserId}})
			if one.Err() != nil {
				if one.Err() == mongo.ErrNoDocuments {
					utils.FailOnError(one.Err(), "No application for those app and user id")
				} else {
					utils.FailOnError(one.Err(), "Failed on searching application from database")
				}
				//TODO: put the message in a dead letter queue
				err = d.Ack(false)
				utils.FailOnError(err, "Failed to acknowledge")
			}
			err = one.Decode(&application)
			utils.FailOnError(err, "Failed to decode application from database")
			//Save the message if raw persistence is activated
			if application.PersistRaw {
				filterCollection := client.Database(cfg.MongoDatabase.Resources).Collection(cfg.MongoCollection.StorageFilters)
				filter, err := utils.CheckFilter(application.RawStorageId, filterCollection)
				if err != nil {
					utils.FailOnError(err, "Failed to get filter")
					err = d.Ack(false)
				}
				if len(filter.Attributes) != 0 {
					filter.Apply(msg)
				}
				persistCollection := client.Database(cfg.MongoDatabase.Data).Collection(application.RawStorageId)
				persistCollection.InsertOne(context.TODO(), msg)
			}
			//Make the acknowledgment
			err = d.Ack(false)
			utils.FailOnError(err, "Failed to acknowledge")
		}
	}()
	fmt.Println("Start consuming...")
	<-forever
	fmt.Println("End of program")
}
