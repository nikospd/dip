package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ApiPort          string `json:"apiPort"`
	SigningKey       string `json:"signingKey"`
	MongoCredentials struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	} `json:"mongoCredentials"`
	MongoDatabase struct {
		Resources string `json:"resources"`
		Data      string `json:"data"`
	} `json:"mongoDatabase"`
	MongoCollection struct {
		Applications string `json:"applications"`
		SourceTokens string `json:"sourceTokens"`
		Storages     string `json:"storages"`
		Users        string `json:"users"`
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
