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

func Connect(URL string, Username string, Password string) (*goftp.Client, error) {

	ftpConfig := goftp.Config{
		User:               Username,
		Password:           Password,
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

func Pull(client *goftp.Client, files []string) ([]string, error) {
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
			// TODO: Q: do we want to use return here?
			return nil, err
		}
		defer buffer.Close()

		if err = client.Retrieve(file, buffer); err != nil {
			log.Error().Msgf("Failed to retreive file with error => %s", err)
			continue
		}
		// TODO: CHECK: if needed becouse list already returns full file path
		local_file := workingDir + "/" + file
		// TODO: Move to seperate function becouse decopress not always needed when pulling e.g. voice station
		if decompressed, err := Decompress(local_file); err != nil {
			log.Error().Msgf("Failed to decompress file '%s' with error => %s", local_file, err)
		} else {
			filePaths = append(filePaths, decompressed)
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

// TODO: Q: maybe rename to ListWithRegex and add List function without regex
func List(client *goftp.Client, path string, pattern string) ([]string, error) {
	path = strings.TrimSuffix(path, "/")
	files, err := client.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files_found []string

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("Failed building regex => %s", err)
	}

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
