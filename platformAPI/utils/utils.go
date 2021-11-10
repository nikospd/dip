package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func CreateRandomHash(len int) string {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func MongoCredentials(user string, password string, host string, port string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
}

func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
	}
}

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	FailOnError(err, "Failed to get hash")
	return string(hash)
}
