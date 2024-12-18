package elevator

import (
	"context"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/decompress"
	"github.com/Kjone1/imageElevator/docker"
	"github.com/rs/zerolog/log"
)

type DockerElevator struct {
	BaseElevator
	registryAdapter docker.RegistryAdapter
	decompressor    decompress.Decompressor
}

const DOCKER_CACHE_FILE = "docker_elevator.json"

func NewDockerElevator(ctx context.Context, baseElevator BaseElevator, registryAdapter docker.RegistryAdapter, elevatorConfig *config.ElevatorConfiguration) *DockerElevator {
	uploadedFiles := loadCache(DOCKER_CACHE_FILE)
	baseElevator.uploadedFiles = uploadedFiles

	var decompressor decompress.Decompressor = new(decompress.TarDecompressor)
	if elevatorConfig.IsUsingXZ {
		decompressor = new(decompress.XZDecompressor)
	}

	elevator := &DockerElevator{
		BaseElevator:    baseElevator,
		registryAdapter: registryAdapter,
		decompressor:    decompressor,
	}

	return elevator
}

func (r *DockerElevator) baseElevator() *BaseElevator {
	return &r.BaseElevator
}

func (r *DockerElevator) Stop() error {
	if err := r.ftpClient.Close(); err != nil {
		log.Error().Msgf("Failed to close connection => %s", err)
		return err
	}
	return nil
}

func (r *DockerElevator) uploadImages() (int, error) {
	tarFiles, err := r.pullFiles()
	if err != nil {
		return 0, err
	}

	tarFiles = r.decompressFiles(tarFiles)

	go func(files []string) {
		for i := range files {
			r.uploadedFiles[filepath.Base(files[i])] = true
		}
		if err := saveCache(DOCKER_CACHE_FILE, r.uploadedFiles); err != nil {
			log.Error().Msgf("Error saving to cache: %s", err)
		}
	}(tarFiles)

	images := tarsToImages(tarFiles)

	fails := 0
	for i := 0; i < len(images); i++ {
		if err := r.registryAdapter.PushTar(r.ctx, &images[i]); err != nil {
			log.Error().Msgf("failed to push %s to registry => %s", images[i].String(), err)
			fails++
			continue
		}

		if err := r.registryAdapter.Sync(r.ctx, &images[i]); err != nil {
			log.Error().Msgf("failed to sync %s:%s => %s", images[i].Name, images[i].Tag, err)
		}
	}

	return len(images) - fails, nil
}

func (r BaseElevator) pullFiles() ([]string, error) {
	remoteFiles, err := r.ftpClient.List(r.workingPath, r.filePattern, r.uploadedFiles)
	if err != nil {
		log.Error().Msgf("Reading FTP directory failed with error => %s", err)
		return nil, err
	}

	if remoteFiles == nil {
		return make([]string, 0), nil
	}

	localFiles, err := r.ftpClient.Pull(remoteFiles...)
	if err != nil {
		log.Error().Msgf("Pulling files from FTP directory => %s", err)
		return nil, err
	}

	return localFiles, nil
}

func (e DockerElevator) decompressFiles(localFiles []string) []string {
	for i, file := range localFiles {
		decompressedFile, err := e.decompressor.Decompress(file)
		if err != nil {
			log.Error().Msgf("failed to decompress %s => %s", file, err)
			continue
		}

		localFiles[i] = decompressedFile
	}
	return localFiles
}

func tarsToImages(tarFiles []string) []docker.Image {
	images := make([]docker.Image, len(tarFiles))
	for i, file := range tarFiles {
		images[i] = *tarToImage(file)
	}
	return images
}

func tarToImage(file string) *docker.Image {
	trimmedFile := filepath.Base(file)
	trimmedFile = strings.TrimSuffix(trimmedFile, ".tar")
	trimmedFile = strings.TrimSuffix(trimmedFile, "-docker")
	trimmedFile = strings.TrimPrefix(trimmedFile, "int-")
	imageName, imageTag := splitTarFile(trimmedFile)
	return &docker.Image{
		Name:    imageName,
		Tag:     imageTag,
		TarPath: file,
	}
}

func splitTarFile(s string) (string, string) {
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		if unicode.IsNumber(r) {
			return s[0 : i-1], s[i:]
		}
	}
	return "", ""
}
