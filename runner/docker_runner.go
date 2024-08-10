package runner

import (
	"context"
	"strings"
	"unicode"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/docker"
	"github.com/Kjone1/imageElevator/ftp"
	"github.com/rs/zerolog/log"
)

type DockerRunner struct {
	RunnerBase
	registryAdapter docker.RegistryAdapter
}

func NewDockerRunner(ctx context.Context, registryAdapter docker.RegistryAdapter, ftpClient ftp.FTPClient, runnerConfig *config.RunnerConfiguration, workingPath, filePattern string) *DockerRunner {

	runner := &DockerRunner{
		RunnerBase:      NewRunnerBase(runnerConfig.SampleRateInMinutes, ftpClient, workingPath, filePattern),
		registryAdapter: registryAdapter,
	}

	return runner
}

func (r *DockerRunner) runnerBase() *RunnerBase {
	return &r.RunnerBase
}

func (r *DockerRunner) Stop() error {
	if err := r.ftpClient.Close(); err != nil {
		log.Error().Msgf("Failed to close connection => %s", err)
		return err
	}
	return nil
}

func (r *DockerRunner) uploadImages() (int, error) {
	tarFiles, err := r.pullFiles()
	if err != nil {
		return 0, err
	}

	images := tarsToImages(tarFiles)

	for i := 0; i < len(images); i++ {
		if err := r.registryAdapter.PushTar(r.ctx, &images[i]); err != nil {
			log.Error().Msgf("failed to push %s to registry => %s", images[i].TarPath, err)
		}
	}

	return len(images), nil
}

func (r RunnerBase) pullFiles() ([]string, error) {
	remoteFiles, err := r.ftpClient.List(r.workingPath, r.filePattern)
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

func tarsToImages(tarFiles []string) []docker.Image {
	images := make([]docker.Image, len(tarFiles))
	for i, file := range tarFiles {
		trimmedFile := strings.TrimSuffix(file, ".tar")
		trimmedFile = strings.TrimSuffix(trimmedFile, "-docker")
		trimmedFile = strings.TrimPrefix(trimmedFile, "int-")
		imageName, imageTag := splitTarFile(trimmedFile)
		image := docker.Image{
			Name:    imageName,
			Tag:     imageTag,
			TarPath: file,
		}
		images[i] = image
	}
	return images
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
