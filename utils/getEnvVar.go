package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Couldn't load .env file")
	}
}

func GetEnvironmentVariable(key string, defaultValue string) string {
	loadEnv()

	value, exists := os.LookupEnv(key)

	if !exists {
		panic(fmt.Sprintf("environment variable %s not found", key))
	}

	if value != "" {
		return value
	}

	return defaultValue
}
