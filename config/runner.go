package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type RunnerConfiguration struct {
	SampleRate time.Duration
	Timer      *time.Timer
	RunnerChan chan interface{}
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
		log.Err(fmt.Errorf("%s is not an integer", sampleRateString)).Msgf("failed to convert sample rate to int")
		sampleRate = 15
	}

	return &RunnerConfiguration{
		SampleRate: time.Duration(sampleRate),
	}
}
