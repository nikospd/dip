package main

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
	fmt.Println("hey")
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
		cfg.AmqpQueues.IntegrationQueue, true, false, false, false, nil)
	utils.FailOnError(err, "Failed to declare integration queue")

	msgs, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg utils.IncomingMessage
			err := json.Unmarshal(d.Body, &msg)
			utils.FailOnError(err, "Failed to read integration message")
			//Search for integrations
			var integrations []utils.Integration
			igrCollection := client.Database(cfg.MongoDatabase.Resources).Collection(cfg.MongoCollection.Integrations)
			cur, err := igrCollection.Find(context.TODO(), bson.D{{"app_id", msg.AppId}})
			cur.All(context.TODO(), &integrations)
			for _, igr := range integrations {

				//Search for existing filters
				filterCollection := client.Database(cfg.MongoDatabase.Resources).Collection(cfg.MongoCollection.StorageFilters)
				filter, _ := utils.CheckFilter(igr.Id, filterCollection)
				if len(filter.Attributes) != 0 {
					filter.Apply(&msg)
				}
				msg.UserId = ""
				msg.AppId = ""
				if err = igr.Send(msg); err != nil {
					utils.FailOnError(err, "Failed to send the message using the integration")
					//TODO: create nack and dead letter queue
				}
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
