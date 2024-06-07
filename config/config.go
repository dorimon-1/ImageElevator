package config

import (
	"fmt"
	"log"
	"os"

	"github.com/containers/image/v5/types"
	"github.com/joho/godotenv"
)

type ContainerConfiguation struct {
	Registry      string
	Repository    string
	SystemContext *types.SystemContext
}

type FtpConfiguration struct {
	FtpServerURL  string
	FtpServerPath string
	FtpUsername   string
	FtpPassword   string
}

var ftpConfig *FtpConfiguration
var containersConfig *ContainerConfiguation

func FtpConfig() FtpConfiguration {
	if ftpConfig == nil {
		ftpConfig = readFtpConfig()
	}

	return *ftpConfig
}

func ContainersConfig() ContainerConfiguation {
	if containersConfig == nil {
		containersConfig = readContainersConfig()
	}

	return *containersConfig
}

func LoadConfig() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("failed reading dotenv: %s", err)
	}

	FtpConfig()
	ContainersConfig()
}

func readFtpConfig() *FtpConfiguration {
	ftpServerURL, err := ReadEnv("FTP_SERVER_URL")
	if err != nil {
		log.Printf("failed to load FTP_SERVER_URL")
	}

	ftpServerPath := ReadEnvWithDefault("FTP_SERVER_PATH", "/")
	ftpUsername := ReadEnvWithDefault("FTP_USERNAME", "ftpuser")
	ftpPassword := ReadEnvWithDefault("FTP_PASSWORD", "ftpuser")

	return &FtpConfiguration{
		FtpServerURL:  ftpServerURL,
		FtpServerPath: ftpServerPath,
		FtpUsername:   ftpUsername,
		FtpPassword:   ftpPassword,
	}
}

func readContainersConfig() *ContainerConfiguation {
	registry, err := ReadEnv("REGISTRY")
	if err != nil {
		log.Fatalf("failed to load REGISTRY")
	}

	repo, err := ReadEnv("REPOSITORY")
	if err != nil {
		log.Fatalf("failed to load REPOSITORY")
	}

	dockerAuthConfig := &types.DockerAuthConfig{
		Username: ReadEnvWithDefault("REPO_USERNAME", "repoUser"),
		Password: ReadEnvWithDefault("REPO_PASSWORD", "repoPass"),
	}

	return &ContainerConfiguation{
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
