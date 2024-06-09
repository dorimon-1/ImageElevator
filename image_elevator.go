package main

import (
	"log"

	"github.com/Kjone1/imageElevator/ftp"
)

func UploadImages() error {
	_, _ = pullTarFile()
	return nil
}

func pullTarFile() ([]string, error) {

	client, err := ftp.Client()
	if err != nil {
		log.Printf("Unable to create FTP client with error => %s", err)
		return nil, err
	}

	images, err := ftp.List(client)
	if err != nil {
		log.Printf("Reading FTP directory failed with error => %s", err)
		return nil, err
	}

	if images == nil {
		log.Printf("No new images were found")
		return nil, err
	}

	files, err := ftp.Pull(client, images)
	if err != nil {
		return nil, err
	}

	return files, nil

}
