package config

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type ElevatorConfiguration struct {
	SampleRateInMinutes time.Duration
	TarRegex            string
	ZipRegex            string
	ZipDestinationPath  string
	IsUsingXZ           bool
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
	zipRegex := ReadEnvWithDefault("ZIP_REGEX", "")
	zipDestinationPath := ReadEnvWithDefault("ZIP_DESTINATION_PATH", "")
	sIsXZ := ReadEnvWithDefault("IS_USING_XZ", "false")
	isZX := sIsXZ == "true"

	return &ElevatorConfiguration{
		SampleRateInMinutes: time.Duration(sampleRate) * time.Minute,
		TarRegex:            tarRegex,
		ZipRegex:            zipRegex,
		ZipDestinationPath:  zipDestinationPath,
		IsUsingXZ:           isZX,
	}
}
