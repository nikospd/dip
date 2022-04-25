package utils

import "fmt"

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
