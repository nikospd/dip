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
	"time"
)

var client *mongo.Client
var channel *amqp.Channel
var queue amqp.Queue
var integrationQueue amqp.Queue
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
		cfg.AmqpQueues.AutomationQueue, true, false, false, false, nil)
	utils.FailOnError(err, "Failed to declare automation queue")
	integrationQueue, err = channel.QueueDeclare(
		cfg.AmqpQueues.IntegrationQueue, true, false, false, false, nil)
	utils.FailOnError(err, "Failed to declare integration queue")

	msgs, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg utils.IncomingMessage
			err := json.Unmarshal(d.Body, &msg)
			utils.FailOnError(err, "Failed to read automation message")
			//Search for automations

			var automations []utils.Automation
			autCollection := client.Database(cfg.MongoDatabase.Resources).Collection("automations")
			cur, err := autCollection.Find(context.TODO(), bson.D{})
			cur.All(context.TODO(), &automations)
			for _, aut := range automations {
				flag, err := aut.Check(msg)
				if err != nil {
					utils.FailOnError(err, "Failed to send the message using the integration")
					//TODO: create nack and dead letter queue
				}
				if flag {
					var integrationMsg utils.IncomingMessage
					integrationMsg.AppId = aut.Id // Using the automation id instead of app id, in order to separate app integrations and automation integrations
					integrationMsg.UserId = aut.UserId
					integrationMsg.ArrivedAt = time.Now()
					integrationMsg.Payload = map[string]interface{}{
						"msg":    "automation activated",
						"reason": fmt.Sprintf("%v %v %v", aut.FirstOperand, aut.Type, aut.SecondOperand),
					}
					msgJs, _ := json.Marshal(integrationMsg)
					err = channel.Publish("", integrationQueue.Name, false, false,
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "text/plain",
							Body:         msgJs,
						})
					utils.FailOnError(err, "Failed to Publish message")
				}
				//Make the acknowledgment
				err = d.Ack(false)
				utils.FailOnError(err, "Failed to acknowledge")
			}
		}
	}()
	fmt.Println("Start consuming...")
	<-forever
	fmt.Println("End of program")
}
