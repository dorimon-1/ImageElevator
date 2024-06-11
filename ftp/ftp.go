package ftp

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ulikunitz/xz"

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

func Pull(client *FtpClient, files []string) ([]string, error){
	filePaths := make([]string, 0)
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

    for _, file := range files {
        log.Printf("Pulling file: %s", file)
        
        buffer, err := os.Create(file)
        if err != nil {
            log.Error().Msgf("Failed to create file with error => %s", err)
            return nil, err
        }
        defer buffer.Close()

        path := fmt.Sprintf("%s/%s", client.FtpServerPath, file)

        if  err = client.Retrieve(path, buffer); err != nil {
            log.Error().Msgf("Failed to retreive file with error => %s", err)
            continue
        }

		file := workingDir + "/" + file
		if filePath, err := Decompress(file); err != nil {
            log.Error().Msgf("Failed to decompress file on path - %s with error => %s", file, err)
        } else {
			filePaths = append(filePaths, filePath)
		}
    } 

    return filePaths, nil
}

func Decompress(inputFilePath string) (string, error) {
    inputFile, err := os.Open(inputFilePath)
    if err != nil {
        return "", err
    }
    defer inputFile.Close()

    outputFilePath := strings.TrimSuffix(inputFilePath, ".xz")
    outputFile, err := os.Create(outputFilePath)
    if err != nil {
        return "", err
    }
    defer outputFile.Close()

    xzReader, err := xz.NewReader(inputFile)
    if err != nil {
        return "", err
    }

    _, err = io.Copy(outputFile, xzReader)
    if err != nil {
        return "", err
    }

    return outputFilePath, nil
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
