package config

import (
	"fmt"
	"log"
	"os"

	"github.com/containers/image/v5/types"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	FtpConfig       *FtpConfig
	ContainerConfig *ContainerConfig
}

type ContainerConfig struct {
	Registry      string
	Repository    string
	SystemContext *types.SystemContext
}

type FtpConfig struct {
	FtpServerURL  string
	FtpServerPath string
	FtpUsername   string
	FtpPassword   string
}

var Config *ServerConfig

func LoadConfig() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("failed reading dotenv: %s", err)
	}

	ftpConfig := readFtpConfig()
	containersConfig := readContainersConfig()

	Config = &ServerConfig{
		FtpConfig:       ftpConfig,
		ContainerConfig: containersConfig,
	}
}

func readFtpConfig() *FtpConfig {
	ftpServerURL, err := ReadEnv("FTP_SERVER_URL")
	if err != nil {
	}

	ftpServerPath := ReadEnvWithDefault("FTP_SERVER_PATH", "/")
	ftpUsername := ReadEnvWithDefault("FTP_USERNAME", "ftpuser")
	ftpPassword := ReadEnvWithDefault("FTP_PASSWORD", "ftpuser")

	return &FtpConfig{
		FtpServerURL:  ftpServerURL,
		FtpServerPath: ftpServerPath,
		FtpUsername:   ftpUsername,
		FtpPassword:   ftpPassword,
	}
}

func readContainersConfig() *ContainerConfig {
	repo, err := ReadEnv("REPOSITORY")
	if err != nil {
	}

	registry, err := ReadEnv("REGISTRY")
	if err != nil {
	}

	dockerAuthConfig := &types.DockerAuthConfig{
		Username: ReadEnvWithDefault("REPO_USERNAME", "repoUser"),
		Password: ReadEnvWithDefault("REPO_PASSWORD", "repoPass"),
	}

	return &ContainerConfig{
		Repository: repo,
		Registry:   registry,
		SystemContext: &types.SystemContext{
			DockerAuthConfig:          dockerAuthConfig,
			DockerCertPath:            ReadEnvWithDefault("DOCKER_CERT_PATH", ""),
			DockerBearerRegistryToken: ReadEnvWithDefault("REGISTRY_BEARER_TOKEN", ""),
		},
	}
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
