package main

import (
	"fmt"
	"log"
	"os"
)

func readEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
		log.Printf("%s environment variable not set, using default: %s", key, defaultValue)
	}

	return value
}

func readEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s environment variable is not defined", key)
	}

	return value, nil
}

type ServerConfig struct {
	FtpServerURL  string
	FtpServerPath string
	FtpUsername   string
	FtpPassword   string
}

var config *ServerConfig

func init() {
	loadConfig()
}

func loadConfig() {
	FtpServerURL, err := readEnv("FTP_SERVER_URL")
	if err != nil {
		panic(err)
	}
	FtpServerPath := readEnvWithDefault("FTP_SERVER_PATH", "/")
	FtpUsername := readEnvWithDefault("FTP_USERNAME", "ftpuser")
	FtpPassword := readEnvWithDefault("FTP_PASSWORD", "ftpuser")

	config = &ServerConfig{
		FtpServerURL:  FtpServerURL,
		FtpServerPath: FtpServerPath,
		FtpUsername:   FtpUsername,
		FtpPassword:   FtpPassword,
	}
}
