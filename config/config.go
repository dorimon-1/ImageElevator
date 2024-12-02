package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	// Use .env file only in dev mode
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(".env"); err != nil {
			log.Error().Msgf("Failed reading dotenv file: %s", err)
		}
	}
}

func LoadConfig() {
	FtpConfig()
	RegistryConfig()
	ElevatorConfig()
}

func ReadIntEnv(key string, defaultValue int) int {
	val, err := strconv.Atoi(ReadEnvWithDefault(key, strconv.Itoa(defaultValue)))
	if err != nil {
		log.Warn().Msgf("Error formatting %s to int => %s", key, err)
		return defaultValue
	}

	return val
}

func ReadEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Warn().Msgf("%s environment variable not set, using default: %s", key, defaultValue)
		return defaultValue
	}

	log.Info().Msgf("Loaded %s=%s", key, value)
	return value
}

func ReadEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s environment variable is not defined", key)
	}

	log.Info().Msgf("Loaded %s=%s", key, value)
	return value, nil
}
