package config

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type RunnerConfiguration struct {
	SampleRateInMinutes time.Duration
}

var runnerConfig *RunnerConfiguration

func RunnerConfig() RunnerConfiguration {
	if runnerConfig == nil {
		runnerConfig = readRunnerConfig()
	}

	return *runnerConfig
}

func readRunnerConfig() *RunnerConfiguration {
	sampleRateString := ReadEnvWithDefault("SAMPLE_RATE_IN_MINUTES", "15")
	sampleRate, err := strconv.Atoi(sampleRateString)
	if err != nil {
		log.Error().Msgf("failed to convert sample rate to int => %s", err)
		sampleRate = 15
	}

	return &RunnerConfiguration{
		SampleRateInMinutes: time.Duration(sampleRate) * time.Minute,
	}
}
