package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secsy/goftp"
)

type FtpClient struct {
	FtpClient *goftp.Client
}

func ftpConnect() (*FtpClient, error) {

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

func ftpPull(client *FtpClient, files []string) {
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", config.FtpServerPath, file)
		log.Printf("Pulling file: %s", path)

		file, err := os.Create(file)
		if err != nil {
			log.Printf("Failed to create file with error => %s", err)
			return
		}
		err = client.FtpClient.Retrieve(path, file)
		if err != nil {
			log.Printf("Failed to retreive file with error => %s", err)
		}
	}
}
func ftpList(client *FtpClient) ([]string, error) {
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
func (client *FtpClient) ftpListEndpoint(c *gin.Context) {
	images, err := ftpList(client)
	if err != nil {
		log.Printf("Reading FTP directory failed with error => %s", err)
		return
	}
	if images == nil {
		log.Printf("No new images were found")
		return
	}
	ftpPull(client, images)
}
