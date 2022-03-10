package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func MongoCredentials(user string, password string, host string, port string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
}

func AmqpCredentials(user string, password string, host string, port string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
}

func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
	}
}

func CheckFilter(storageId string, collection *mongo.Collection) (StorageFilter, error) {
	one := collection.FindOne(context.TODO(), bson.D{{"storage_id", storageId}})
	var filter StorageFilter
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return StorageFilter{}, nil
		}
		return StorageFilter{}, one.Err()
	}
	err := one.Decode(&filter)
	if err != nil {
		return StorageFilter{}, err
	}
	return filter, nil
}
