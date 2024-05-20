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
		log.Fatal("Error loading .env file")
	}
}

func GetEnvironmentVariable(key string, defaultValue string) (string, error) {
	loadEnv()

	value, error := os.LookupEnv(key)

	if error {
		return value, fmt.Errorf("environment variable %s not found", key)
	}

	if value != "" {
		return value, nil
	}

	return defaultValue, nil
}
