package env

import (
	"log"
	"os"
)

func GetWithDefault(key, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVal
}

func GetOrFail(key string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		log.Fatalf("%s environment variable is not set but required", key)
	}
	return value
}
