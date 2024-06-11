package runner

import (
	"context"
	"errors"
	"time"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/containers"
	"github.com/Kjone1/imageElevator/ftp"
	"github.com/rs/zerolog/log"
)

type Runner struct {
	ctx            context.Context
	sampleRate     time.Duration
	timer          *time.Timer
	runUploadChan  chan interface{}
	resetTimerChan chan interface{}
}

func NewRunner(ctx context.Context) *Runner {
	runnerConf := config.RunnerConfig()
	rate := runnerConf.SampleRate * time.Second
	timer := time.NewTimer(rate)
	runUploadChan := make(chan interface{})
	resetTimerChan := make(chan interface{})

	runner := &Runner{ctx: ctx, sampleRate: rate, timer: timer, runUploadChan: runUploadChan, resetTimerChan: resetTimerChan}

	go runner.timerRoutine()
	go runner.uploaderRoutine()

	return runner
}

func (r *Runner) TriggerUpload() {
	r.runUploadChan <- nil
}

func (r *Runner) uploaderRoutine() {
	log.Debug().Msg("Image uploader routine started")
	for {
		select {
		case <-r.runUploadChan:
			if err := r.uploadImages(); err != nil {
				if err.Error() == "no new images were found" {
					log.Debug().Msg(err.Error())
				} else {
					log.Err(err).Msgf("failed to upload images from Image Uploader")
				}
			}
			r.resetTimerChan <- nil

		case <-r.ctx.Done():
			if err := ftp.Close(); err != nil {
				log.Err(err).Msg("failed to close connection")
			}
			close(r.runUploadChan)
			close(r.resetTimerChan)
			return
		}
	}
}

func (r *Runner) uploadImages() error {
	tarFiles, err := pullTarFile()
	if err != nil {
		return err
	}

	if err := containers.PushMultipleTars(r.ctx, tarFiles, "imageName", "tag", config.ContainersConfig()); err != nil {
		return err
	}

	return nil
}

func pullTarFile() ([]string, error) {

	client, err := ftp.Client()
	if err != nil {
		log.Err(err).Msg("Unable to create FTP client with error")
		return nil, err
	}

	images, err := ftp.List(client)
	if err != nil {
		log.Err(err).Msg("Reading FTP directory failed with error")
		return nil, err
	}

	if images == nil {
		return nil, errors.New("no new images were found")
	}

	files, err := ftp.Pull(client, images)
	if err != nil {
		log.Err(err).Msg("Pulling files from FTP directory")
		return nil, err
	}

	return files, nil
}
