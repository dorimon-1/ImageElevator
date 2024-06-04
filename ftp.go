package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secsy/goftp"
)

func ftpConnect() (*goftp.Client, error) {

	ftpConfig := goftp.Config{
		User:               config.FtpUsername,
		Password:           config.FtpPassword,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stdout,
	}

	client, err := goftp.DialConfig(ftpConfig, config.FtpServerURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ftpListEndpoint(c *gin.Context) {
	client, err := ftpConnect()
	if err != nil {
		log.Printf("Failed to create FTP client with error => %s", err)
		return
	}
	files, err := client.ReadDir(config.FtpServerPath)
	if err != nil {
		log.Printf("Reading FTP directory failed with error => %s", err)
		return
	}

	for _, file := range files {
		pattern := "^int-"
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			log.Printf("Failed while matching %s againts regex pattern '%s' with error => %s", file.Name(), pattern, err)
		}

		if matched {
			log.Printf("Found file: %s", file.Name())
		}
	}
}
