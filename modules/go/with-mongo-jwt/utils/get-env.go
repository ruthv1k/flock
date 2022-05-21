package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Gets a specified environment variable
//
// GetEnv(variableName)
func GetEnv(v string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return "", err
	}

	e := os.Getenv(v)

	return e, nil
}
