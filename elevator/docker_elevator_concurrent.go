package elevator

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/decompress"
	"github.com/Kjone1/imageElevator/docker"
	"github.com/rs/zerolog/log"
)

type ConcurrentDockerElevator struct {
	BaseElevator
	registryAdapter docker.RegistryAdapter
	decompressor    decompress.Decompressor
}

func NewConcurrentDockerElevator(ctx context.Context, baseElevator BaseElevator, registryAdapter docker.RegistryAdapter, elevatorConfig *config.ElevatorConfiguration) *ConcurrentDockerElevator {
	uploadedFiles := loadCache(DOCKER_CACHE_FILE)
	baseElevator.uploadedFiles = uploadedFiles

	var decompressor decompress.Decompressor = new(decompress.TarDecompressor)
	if elevatorConfig.IsUsingXZ {
		decompressor = new(decompress.XZDecompressor)
	}

	elevator := &ConcurrentDockerElevator{
		BaseElevator:    baseElevator,
		registryAdapter: registryAdapter,
		decompressor:    decompressor,
	}

	return elevator
}

func (r *ConcurrentDockerElevator) baseElevator() *BaseElevator {
	return &r.BaseElevator
}

func (r *ConcurrentDockerElevator) Stop() error {
	if err := r.ftpClient.Close(); err != nil {
		log.Error().Msgf("Failed to close connection => %s", err)
		return err
	}
	return nil
}

func (r *ConcurrentDockerElevator) uploadImages() (int, error) {
	files, err := r.pullFiles()
	if err != nil {
		return 0, err
	}
	uploadedFilesChan := make(chan string)
	for _, file := range files {
		go func(file string, filesChan chan string) {
			outputFile, err := r.uploadFile(file)
			if err != nil {
				log.Error().Msg(err.Error())
			}
			uploadedFilesChan <- outputFile
		}(file, uploadedFilesChan)
	}

	uploads := 0
	for range files {
		file := <-uploadedFilesChan
		if file == "" {
			continue
		}
		r.uploadedFiles[file] = true
		uploads++

	}

	close(uploadedFilesChan)

	if err := saveCache(DOCKER_CACHE_FILE, r.uploadedFiles); err != nil {
		log.Error().Msgf("Error saving to cache: %s", err)
	}

	return uploads, nil
}

func (r *ConcurrentDockerElevator) uploadFile(file string) (string, error) {
	decompressedFile, err := r.decompressor.Decompress(file)
	if err != nil {
		log.Error().Msgf("failed to decompress file %s => %s", file, err)
		return "", err
	}

	if err := os.Remove(file); err != nil {
		log.Error().Msgf("failed to remove file %s => %s", file, err)
	}

	image := tarToImage(decompressedFile)

	if err := r.registryAdapter.PushTar(r.ctx, image); err != nil {
		log.Error().Msgf("failed to push %s to registry => %s", image.String(), err)
		return "", err
	}

	if err := r.registryAdapter.Sync(r.ctx, image); err != nil {
		log.Error().Msgf("failed to sync %s:%s => %s", image.Name, image.Tag, err)
	}

	if err := os.Remove(decompressedFile); err != nil {
		log.Error().Msgf("failed to remove file %s => %s", file, err)
	}

	return filepath.Base(file), nil
}
