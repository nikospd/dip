package main

import (
	"context"
	"dev.com/utils"
	"fmt"
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
		Connect to MongoDB
	*/
	clientOptions := options.Client().ApplyURI("mongodb://test:test@localhost:27017")
	client, connectionError := mongo.Connect(context.TODO(), clientOptions)
	if connectionError != nil {
		log.Fatalln(connectionError)
	}
	connectionError = client.Ping(context.TODO(), nil)
	if connectionError != nil {
		log.Fatalln(connectionError)
	}
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
				go v.ExecuteTask()
			}
		} else {
			fmt.Println("No remaining tasks")
		}
		time.Sleep(time.Duration(searchingTasksInterval) * time.Second)
	}
}

func createTasks(numberOfTasks int) []utils.PullSourceTask {
	var sourceTable []utils.PullSourceTask
	for i := 0; i < numberOfTasks+0; i++ {
		interval := rand.Intn(20)
		//d := time.Duration(interval) * time.Second
		source := new(utils.PullSourceTask)
		source.TaskId = fmt.Sprintf("%v", i)
		source.Interval = interval
		source.NextExecution = time.Now()
		sourceTable = append(sourceTable, *source)
	}
	return sourceTable
}
