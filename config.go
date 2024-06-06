package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	FtpServerURL  string
	FtpServerPath string
	FtpUsername   string
	FtpPassword   string
}

var config *ServerConfig

func LoadConfig() {
	godotenv.Load()

	FtpServerURL, err := ReadEnv("FTP_SERVER_URL")
	if err != nil {
		panic(err)
	}
	FtpServerPath := ReadEnvWithDefault("FTP_SERVER_PATH", "/")
	FtpUsername := ReadEnvWithDefault("FTP_USERNAME", "ftpuser")
	FtpPassword := ReadEnvWithDefault("FTP_PASSWORD", "ftpuser")

	config = &ServerConfig{
		FtpServerURL:  FtpServerURL,
		FtpServerPath: FtpServerPath,
		FtpUsername:   FtpUsername,
		FtpPassword:   FtpPassword,
	}

	fmt.Printf("** Server is loaded with the following: %s ** \n", *config)
}

func ReadEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
		log.Printf("%s environment variable not set, using default: %s", key, defaultValue)
	}

	return value
}

func ReadEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s environment variable is not defined", key)
	}

	return value, nil
}

