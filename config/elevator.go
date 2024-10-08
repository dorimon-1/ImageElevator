package config

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type ElevatorConfiguration struct {
	SampleRateInMinutes time.Duration
	TarRegex            string
}

var elevatorConfig *ElevatorConfiguration

func ElevatorConfig() ElevatorConfiguration {
	if elevatorConfig == nil {
		elevatorConfig = readElevatorConfig()
	}

	return *elevatorConfig
}

func readElevatorConfig() *ElevatorConfiguration {
	sampleRateString := ReadEnvWithDefault("SAMPLE_RATE_IN_MINUTES", "15")
	sampleRate, err := strconv.Atoi(sampleRateString)
	if err != nil {
		log.Error().Msgf("failed to convert sample rate to int => %s", err)
		sampleRate = 15
	}

	tarRegex := ReadEnvWithDefault("TAR_REGEX", "")

	return &ElevatorConfiguration{
		SampleRateInMinutes: time.Duration(sampleRate) * time.Minute,
		TarRegex:            tarRegex,
	}
}
