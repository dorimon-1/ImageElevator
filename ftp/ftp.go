package ftp

import (
	"fmt"
	"github.com/Kjone1/imageElevator/config"
	"github.com/rs/zerolog/log"
	"github.com/secsy/goftp"
	"github.com/ulikunitz/xz"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var client *goftp.Client

func Client() (*goftp.Client, error) {
	if client == nil {
		config := config.FtpConfig()
		ftpClient, err := Connect(config.FtpServerURL, config.FtpUsername, config.FtpPassword)
		if err != nil {
			return nil, err
		}
		client = ftpClient
	}
	return client, nil
}

func Connect(URL string, username string, password string) (*goftp.Client, error) {

	ftpConfig := goftp.Config{
		User:               username,
		Password:           password,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             nil,
	}

	client, err := goftp.DialConfig(ftpConfig, URL)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Pull(client *goftp.Client, files ...string) ([]string, error) {
	filesPulled := make([]string, 0)
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	for _, remote_file := range files {
		local_file := fmt.Sprintf("%s/%s", workingDir, filepath.Base(remote_file))
		log.Info().Msgf("Pulling file '%s' from remote to %s", remote_file, local_file)

		buffer, err := os.Create(local_file)
		if err != nil {
			log.Error().Msgf("Failed to create file with error => %s", err)
			continue
		}
		defer buffer.Close()

		if err = client.Retrieve(remote_file, buffer); err != nil {
			log.Error().Msgf("Failed to retreive file with error => %s", err)
			continue
		}
		filesPulled = append(filesPulled, local_file)

	}

	return filesPulled, nil
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

func List(client *goftp.Client, path string, pattern string) ([]string, error) {

	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}

	files, err := client.ReadDir(path)
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed building regex => %s", err)
	}

	var files_found []string

	for _, file := range files {
		matched := regex.MatchString(file.Name())

		if matched {
			log.Info().Msgf("Found file: %s", file.Name())
			var full_file_path string
			if path == "/" {
				full_file_path = fmt.Sprintf("/%s", file.Name())
			} else {
				full_file_path = fmt.Sprintf("%s/%s", path, file.Name())
			}
			files_found = append(files_found, full_file_path)
		}
	}
	return files_found, nil
}
