package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ApiPort          string `json:"apiPort"`
	MongoCredentials struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	} `json:"mongoCredentials"`
	AmqpCredentials struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	} `json:"amqpCredentials"`
	AmqpQueues struct {
		IncomingData string `json:"incomingData"`
	} `json:"amqpQueues"`
	MongoDatabase struct {
		Resources string `json:"resources"`
	} `json:"mongoDatabase"`
	MongoCollection struct {
		PullSources string `json:"pullSources"`
	} `json:"mongoCollection"`
}

func ReadConf(configFile string, cfg *Configuration) {
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err, "Failed to read configuration file")
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatalln(err, "Failed to read configuration file")
	}
}
