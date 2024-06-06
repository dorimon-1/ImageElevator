package ftp

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/secsy/goftp"
)

type FtpClient struct {
	FtpClient *goftp.Client
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
	return &FtpClient{FtpClient: client}, nil
}

func Pull(client *FtpClient, files []string) {
	for _, file := range files {

		log.Printf("Pulling file: %s", file)

		buffer, err := os.Create(file)
		if err != nil {
			log.Printf("Failed to create file with error => %s", err)
			return
		}
		path := fmt.Sprintf("%s/%s", config.FtpServerPath, file)
		err = client.FtpClient.Retrieve(path, buffer)
		if err != nil {
			log.Printf("Failed to retreive file with error => %s", err)
		}
	}
}
func List(client *FtpClient) ([]string, error) {
	files, err := client.FtpClient.ReadDir(config.FtpServerPath)
	if err != nil {
		return nil, err
	}
	var files_found []string
	for _, file := range files {
		pattern := "^int-.*-docker$"
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			log.Printf("Failed while matching %s againts regex pattern '%s' with error => %s", file.Name(), pattern, err)
		}

		if matched {
			log.Printf("Found file: %s", file.Name())
			files_found = append(files_found, file.Name())
		}
	}
	return files_found, nil
}
