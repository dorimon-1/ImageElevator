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
	syncRegistriesTokens := ReadEnvWithDefault("SYNC_REGISTRIES_BEARER_TOKEN", "")

	if syncRegistries == "" {
		return nil
	}

	registries := strings.Split(syncRegistries, ",")
	tokens := strings.Split(syncRegistriesTokens, ",")
	if len(tokens) != len(registries) {
		log.Error().Msgf("failed to load sync registries, missing token or registry: %d registries, %d tokens", len(registries), len(tokens))
	}

	registriesConfig := make([]RegistryConfiguration, len(registries))

	for i, registry := range registries {
		token := tokens[i]
		registriesConfig[i] = RegistryConfiguration{
			Registry:   registry,
			Repository: RegistryConfig().Repository,
			SystemContext: &types.SystemContext{
				DockerBearerRegistryToken: token,
			},
		}
	}

	return registriesConfig
}
