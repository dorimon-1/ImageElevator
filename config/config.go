package config

import (
	"fmt"
	"log"
	"os"

	"github.com/containers/image/v5/types"
)

type ServerConfig struct {
	*FtpConfig
	*ContainerConfig
}

type ContainerConfig struct {
	Registry      string
	RepoURL       string
	SystemContext *types.SystemContext
}

type FtpConfig struct {
	FtpServerURL  string
	FtpServerPath string
	FtpUsername   string
	FtpPassword   string
}

var config *ServerConfig

func init() {
	config = readConfig()
}

func Config() *ServerConfig {
	return config
}

func readConfig() *ServerConfig {
	ftpConfig := readFtpConfig()
	containersConfig := readContainersConfig()

	return &ServerConfig{
		FtpConfig:       ftpConfig,
		ContainerConfig: containersConfig,
	}
}

func readFtpConfig() *FtpConfig {
	ftpServerURL, err := readEnv("FTP_SERVER_URL")
	if err != nil {
	}

	ftpServerPath := readEnvWithDefault("FTP_SERVER_PATH", "/")
	ftpUsername := readEnvWithDefault("FTP_USERNAME", "ftpuser")
	ftpPassword := readEnvWithDefault("FTP_PASSWORD", "ftpuser")

	return &FtpConfig{
		FtpServerURL:  ftpServerURL,
		FtpServerPath: ftpServerPath,
		FtpUsername:   ftpUsername,
		FtpPassword:   ftpPassword,
	}
}

func readContainersConfig() *ContainerConfig {
	repoURL, err := readEnv("REPO_URL")
	if err != nil {
	}

	registry, err := readEnv("REGISTRY")
	if err != nil {
	}

	dockerAuthConfig := &types.DockerAuthConfig{
		Username: readEnvWithDefault("REPO_USERNAME", "admin"),
		Password: readEnvWithDefault("REPO_PASSWORD", "admin"),
	}

	return &ContainerConfig{
		RepoURL:  repoURL,
		Registry: registry,
		SystemContext: &types.SystemContext{
			DockerAuthConfig:          dockerAuthConfig,
			DockerCertPath:            readEnvWithDefault("DOCKER_CERT_PATH", ""),
			DockerBearerRegistryToken: readEnvWithDefault("REGISTRY_BEARER_TOKEN", ""),
		},
	}
}

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
