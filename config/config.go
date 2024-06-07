package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("failed reading dotenv: %s", err)
	}
}

func LoadConfig() {
	FtpConfig()
	ContainersConfig()
}

func ReadEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("%s environment variable not set, using default: %s", key, defaultValue)
		return defaultValue
	}

	log.Printf("Loaded %s=%s", key, value)
	return value
}

func ReadEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s environment variable is not defined", key)
	}

	log.Printf("Loaded %s=%s", key, value)
	return value, nil
}
