package ftp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Kjone1/imageElevator/config"
	"github.com/prasad83/goftp"
	"github.com/rs/zerolog/log"
)

type GoFTP struct {
	client *goftp.Client
}

var ftpClient *GoFTP

func Client() (*GoFTP, error) {
	if ftpClient == nil {
		config := config.FtpConfig()
		goFTPClient, err := Connect(config.FtpServerURL, config.FtpUsername, config.FtpPassword, os.Stdout)
		if err != nil {
			return nil, err
		}
		_ = goFTPClient

		ftpClient = &GoFTP{
			client: goFTPClient,
		}
	}
	return ftpClient, nil
}

func Connect(URL string, username string, password string, logger io.Writer) (*goftp.Client, error) {
	ftpConfig := goftp.Config{
		User:               username,
		Password:           password,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             logger,
	}

	client, err := goftp.DialConfig(ftpConfig, URL)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (f *GoFTP) Close() error {
	if f.client != nil {
		return f.client.Close()
	}
	return errors.New("connection is already closed")
}

func (f *GoFTP) Pull(files ...string) ([]string, error) {
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

		if err = f.client.Retrieve(remote_file, buffer); err != nil {
			log.Error().Msgf("Failed to retreive file with error => %s", err)
			continue
		}
		filesPulled = append(filesPulled, local_file)

	}

	return filesPulled, nil
}

func (f *GoFTP) List(path string, pattern string, bannedFiles map[string]bool) ([]string, error) {
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}

	files, err := f.client.ReadDir(path)
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

		if matched && !isBanned(bannedFiles, file.Name()) {
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

func isBanned(bannedFiles map[string]bool, file string) bool {
	return bannedFiles[file]
}
