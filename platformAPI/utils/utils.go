package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
