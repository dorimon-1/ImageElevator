package config

import (
	"strings"

	"github.com/containers/image/v5/types"
	"github.com/rs/zerolog/log"
)

var syncConfig []RegistryConfiguration

func SyncConfig() []RegistryConfiguration {
	if syncConfig == nil {
		syncConfig = readSyncConfig()
	}

	return syncConfig
}

func readSyncConfig() []RegistryConfiguration {
	syncRegistries := ReadEnvWithDefault("SYNC_REGISTRIES", "")
	syncRepositories := ReadEnvWithDefault("SYNC_REPOSITORIES", "")
	syncRegistriesTokens := ReadEnvWithDefault("SYNC_REGISTRIES_BEARER_TOKEN", "")

	if syncRegistries == "" || syncRepositories == "" {
		return nil
	}

	registries := strings.Split(syncRegistries, ",")
	repositories := strings.Split(syncRepositories, ",")
	tokens := strings.Split(syncRegistriesTokens, ",")
	if len(tokens) != len(registries) || len(registries) != len(repositories) {
		log.Error().Msgf("failed to load sync registries, missing token, repository or registry: %d registries, %d tokens, %d repositories", len(registries), len(tokens), len(repositories))
	}

	registriesConfig := make([]RegistryConfiguration, len(registries))

	for i, registry := range registries {
		token := tokens[i]
		registriesConfig[i] = RegistryConfiguration{
			Registry:   registry,
			Repository: registries[i],
			SystemContext: &types.SystemContext{
				DockerBearerRegistryToken:   token,
				DockerInsecureSkipTLSVerify: types.OptionalBoolTrue,
			},
		}
	}

	return registriesConfig
}
