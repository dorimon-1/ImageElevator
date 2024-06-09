package ftp

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"time"

	"github.com/Kjone1/imageElevator/config"
	"github.com/secsy/goftp"
)

type FtpClient struct {
	*goftp.Client
	*config.FtpConfiguration
}

var client *FtpClient

func Client() (*FtpClient, error) {
	if client == nil {
		ftpClient, err := Connect()
		if err != nil {
			return nil, err
		}
		client = ftpClient
	}
	return client, nil
}

func Connect() (*FtpClient, error) {
	config := config.FtpConfig()

	ftpConfig := goftp.Config{
		User:               config.FtpUsername,
		Password:           config.FtpPassword,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             nil,
	}

	client, err := goftp.DialConfig(ftpConfig, config.FtpServerURL)
	if err != nil {
		return nil, err
	}

	return &FtpClient{Client: client, FtpConfiguration: &config}, nil
}

func Pull(client *FtpClient, files []string) {
	for _, file := range files {

		log.Printf("Pulling file: %s", file)

		buffer, err := os.Create(file)
		if err != nil {
			log.Error().Msgf("Failed to create file with error => %s", err)
			return
		}
		path := fmt.Sprintf("%s/%s", client.FtpServerPath, file)
		err = client.Retrieve(path, buffer)
		if err != nil {
			log.Error().Msgf("Failed to retreive file with error => %s", err)
		}
	}
}

func List(client *FtpClient) ([]string, error) {
	files, err := client.ReadDir(client.FtpServerPath)
	if err != nil {
		return nil, err
	}

	var files_found []string

	//TODO: make pattern an environement variable
	pattern := "^int-.*-docker$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("building regex: %s", err)
	}

	for _, file := range files {
		matched := regex.MatchString(file.Name())

		if matched {
			log.Info().Msgf("Found file: %s", file.Name())
			files_found = append(files_found, file.Name())
		}
	}
	return files_found, nil
}
