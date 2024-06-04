package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secsy/goftp"
)

func ftpConnect() *goftp.Client {

	ftpConfig := goftp.Config{
		User:               config.FtpUsername,
		Password:           config.FtpPassword,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stdout,
	}

	client, err := goftp.DialConfig(ftpConfig, config.FtpServerURL)
	if err != nil {
		panic(err)
	}
	return client
}

func ftpListEndpoint(c *gin.Context) {
	client := ftpConnect()

	files, err := client.ReadDir(config.FtpServerPath)
	if err != nil {
		log.Printf("Reading FTP directory failed with error => %s", err)
	}

	for _, file := range files {
		matched, err := regexp.MatchString("^int-", file.Name())
		if err != nil {
			log.Printf("Failed to match file name againts regex with error => %s", err)
		}

		if matched {
			log.Printf("Found file: %s", file.Name())
		}
	}
}
