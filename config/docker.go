package config

import (
	"github.com/rs/zerolog/log"

	"github.com/containers/image/v5/types"
)

type RegistryConfiguration struct {
	Registry      string
	Repository    string
	SystemContext *types.SystemContext
}

var containersConfig *RegistryConfiguration

func RegistryConfig() RegistryConfiguration {
	if containersConfig == nil {
		containersConfig = readRegistryConfig()
	}

	return *containersConfig
}

func readRegistryConfig() *RegistryConfiguration {
	registry, err := ReadEnv("REGISTRY")
	if err != nil {
		log.Fatal().Msg("Failed to load REGISTRY env var")
	}

	repo, err := ReadEnv("REPOSITORY")
	if err != nil {
		log.Fatal().Msg("Failed to load REPOSITORY env var")
	}

	dockerAuthConfig := &types.DockerAuthConfig{
		Username: ReadEnvWithDefault("REPO_USERNAME", "repoUser"),
		Password: ReadEnvWithDefault("REPO_PASSWORD", "repoPass"),
	}

	return &RegistryConfiguration{
		Repository: repo,
		Registry:   registry,
		SystemContext: &types.SystemContext{
			DockerAuthConfig:            dockerAuthConfig,
			DockerCertPath:              ReadEnvWithDefault("DOCKER_CERT_PATH", ""),
			DockerBearerRegistryToken:   ReadEnvWithDefault("REGISTRY_BEARER_TOKEN", ""),
			DockerInsecureSkipTLSVerify: types.OptionalBoolTrue,
		},
	}
}
