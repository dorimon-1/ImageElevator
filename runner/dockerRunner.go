package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/docker"
	"github.com/Kjone1/imageElevator/ftp"
	"github.com/rs/zerolog/log"
)

type DockerRunner struct {
	BasicRunner
	registryAdapter docker.RegistryAdapter
	ftpClient       ftp.FTPClient
	workingPath     string
	filePattern     string
}

func NewDockerRunner(ctx context.Context, registryAdapter docker.RegistryAdapter, ftpClient ftp.FTPClient, runnerConfig *config.RunnerConfiguration, workingPath, filePattern string) *DockerRunner {
	timer := time.NewTimer(runnerConfig.SampleRateInMinutes)
	runUploadChan := make(chan any, 1)
	resetTimerChan := make(chan any)

	runner := &DockerRunner{
		BasicRunner: BasicRunner{
			ctx:            ctx,
			sampleRate:     runnerConfig.SampleRateInMinutes,
			timer:          timer,
			runUploadChan:  runUploadChan,
			resetTimerChan: resetTimerChan,
		},
		ftpClient:       ftpClient,
		registryAdapter: registryAdapter,
		workingPath:     workingPath,
		filePattern:     filePattern,
	}

	return runner
}

func (r *DockerRunner) basicRunner() *BasicRunner {
	return &r.BasicRunner
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

func (r *DockerRunner) pullFiles() ([]string, error) {
	fmt.Println("Uploading..")
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

// TODO: Make a function receives a list of tar files and returns a docker.Image (ImageName, Tag, TarPath) by regex
func tarsToImages(tarFiles []string) []docker.Image {
	_ = tarFiles
	return nil
}
