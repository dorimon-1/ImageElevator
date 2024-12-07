package elevator

import (
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"path/filepath"
)

type ZipElevator struct {
	BaseElevator
	destinationPath string
}

const ZIP_CACHE_FILE = "zip_elevator.json"

func NewZipElevator(ctx context.Context, destPath string, baseElevator BaseElevator) *ZipElevator {
	return &ZipElevator{
		BaseElevator:    baseElevator,
		destinationPath: destPath,
	}
}

func (r *ZipElevator) baseElevator() *BaseElevator {
	return &r.BaseElevator
}
func (r *ZipElevator) Stop() error {
	return nil
}
func (r *ZipElevator) uploadImages() (int, error) {
	count := 0
	zipFiles, err := r.pullFiles()
	if err != nil {
		return 0, err
	}

	go func(files []string) {
		for i := range files {
			r.uploadedFiles[filepath.Base(files[i])] = true
		}
		if err := saveCache(ZIP_CACHE_FILE, r.uploadedFiles); err != nil {
			log.Error().Msgf("Error saving to cache: %s", err)
		}
	}(zipFiles)

	for _, file := range zipFiles {
		fileName := path.Base(file)
		if fileName == "/" || fileName == "." {
			log.Error().Msg("Couldn't find file")
			continue
		}
		dest := r.destinationPath + "/" + fileName
		if err := copyFiles(file, dest); err != nil {
			log.Error().Err(err)
			continue
		}

		count++
	}

	return count, nil
}

// NOTE: destinationPath should be the path including the file name!!!
func copyFiles(file string, destinationPath string) error {
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	return nil
}
