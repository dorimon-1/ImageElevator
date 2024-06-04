package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secsy/goftp"
)

func ftpListEndpoint(c *gin.Context) {
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

	files, err := client.ReadDir(config.FtpServerPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		matched, err := regexp.MatchString("^int-", file.Name())
		if err != nil {
			panic(err)
		}

		if matched {
			log.Printf("Found file: %s", file.Name())
		}

	}
}

func healthEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
