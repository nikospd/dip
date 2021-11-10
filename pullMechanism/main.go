package main

import (
	"context"
	"dev.com/config"
	"dev.com/utils"
	"fmt"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("start of the program")
	rand.Seed(time.Now().UnixNano())
	/*
		Read configuration file
	*/
	var cfg config.Configuration
	config.ReadConf("config.json", &cfg)
	/*
		/*
			Connect to MongoDB
	*/
	mongoCred := cfg.MongoCredentials
	mongoUri := utils.MongoCredentials(mongoCred.User, mongoCred.Password, mongoCred.Host, mongoCred.Port)
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, connectionError := mongo.Connect(context.TODO(), clientOptions)
	if connectionError != nil {
		log.Fatalln(connectionError)
	}
	connectionError = client.Ping(context.TODO(), nil)
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
	channel, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	queue, err := channel.QueueDeclare(
		cfg.AmqpQueues.IncomingData, true, false, false, false, nil)
	utils.FailOnError(err, "Failed to declare a queue")
	/*
		Get N last tasks shorted by nextExecution and execute them in different goroutines
		every X seconds
	*/
	searchingTasksInterval := 2
	col := client.Database("staging").Collection("pull_sources")
	var tasks []utils.PullSourceTask
	for true {
		cur, err := col.Find(context.TODO(), bson.D{{"next_execution", bson.D{{"$lt", time.Now()}}}})
		if err != nil {
			fmt.Println(err)
		} else if cur.RemainingBatchLength() != 0 {
			cur.All(context.TODO(), &tasks)
			fmt.Println("found tasks now execute them")
			for _, v := range tasks {
				/*
					Update next Execution time
				*/
				v.HandleTask(col)
				/*
					Actual execute the task
				*/
				go v.ExecuteTask(channel, queue)
			}
		} else {
			fmt.Println("No remaining tasks")
		}
		time.Sleep(time.Duration(searchingTasksInterval) * time.Second)
	}
}
