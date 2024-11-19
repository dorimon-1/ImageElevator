package config

import "github.com/rs/zerolog/log"

type FtpConfiguration struct {
	FtpServerURL     string
	FtpServerPath    string
	FtpUsername      string
	FtpPassword      string
	FtpLoggerEnabled bool
}

var ftpConfig *FtpConfiguration

func FtpConfig() FtpConfiguration {
	if ftpConfig == nil {
		ftpConfig = readFtpConfig()
	}

	return *ftpConfig
}

func readFtpConfig() *FtpConfiguration {
	ftpServerURL, err := ReadEnv("FTP_SERVER_URL")
	if err != nil {
		log.Error().Msg("Failed to load FTP_SERVER_URL env var")
	}

	ftpServerPath := ReadEnvWithDefault("FTP_SERVER_PATH", "/")
	ftpUsername := ReadEnvWithDefault("FTP_USERNAME", "ftpuser")
	ftpPassword := ReadEnvWithDefault("FTP_PASSWORD", "ftpuser")
	ftpLoggerEnabled := ReadEnvWithDefault("FTP_LOGGER_ENABLED", "")
	loggerEnabled := false
	if ftpLoggerEnabled != "" {
		loggerEnabled = true
	}

	return &FtpConfiguration{
		FtpServerURL:     ftpServerURL,
		FtpServerPath:    ftpServerPath,
		FtpUsername:      ftpUsername,
		FtpPassword:      ftpPassword,
		FtpLoggerEnabled: loggerEnabled,
	}
}
