package runner

import (
	"context"
	"time"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/docker"
	"github.com/Kjone1/imageElevator/ftp"
	"github.com/rs/zerolog/log"
)

type DockerRunner struct {
	ctx             context.Context
	sampleRate      time.Duration
	timer           *time.Timer
	runUploadChan   chan interface{}
	resetTimerChan  chan interface{}
	registryAdapter docker.RegistryAdapter
	ftpClient       ftp.FTPClient
	workingPath     string
	filePattern     string
}

func NewDockerRunner(ctx context.Context, registryAdapter docker.RegistryAdapter, ftpClient ftp.FTPClient, runnerConfig *config.RunnerConfiguration) *DockerRunner {
	timer := time.NewTimer(runnerConfig.SampleRateInMinutes)
	runUploadChan := make(chan interface{}, 1)
	resetTimerChan := make(chan interface{})

	runner := &DockerRunner{
		ctx:             ctx,
		sampleRate:      runnerConfig.SampleRateInMinutes,
		timer:           timer,
		runUploadChan:   runUploadChan,
		resetTimerChan:  resetTimerChan,
		registryAdapter: registryAdapter,
	}

	return runner
}

func (r *DockerRunner) Start() {
	go r.timerRoutine()
	go uploaderRoutine(r)
}

func (r *DockerRunner) uploadImages() (int, error) {
	tarFiles, err := r.pullFiles()
	if err != nil {
		return 0, err
	}

	images := tarsToImages(tarFiles)

	for i := 0; i < len(images); i++ {
		if err := r.registryAdapter.PushTar(r.ctx, &images[i]); err != nil {
			log.Err(err).Msgf("failed to push %s to registry", images[i].TarPath)
		}
	}

	return len(images), nil
}

func (r *DockerRunner) pullFiles() ([]string, error) {
	remoteFiles, err := r.ftpClient.List(r.workingPath, r.filePattern)
	if err != nil {
		log.Err(err).Msg("Reading FTP directory failed with error")
		return nil, err
	}

	if remoteFiles == nil {
		return make([]string, 0), nil
	}

	localFiles, err := r.ftpClient.Pull(remoteFiles...)
	if err != nil {
		log.Err(err).Msg("Pulling files from FTP directory")
		return nil, err
	}

	return localFiles, nil
}

// TODO: Make a function receives a list of tar files and returns a docker.Image (ImageName, Tag, TarPath) by regex
func tarsToImages(tarFiles []string) []docker.Image {
	_ = tarFiles
	return nil
}
