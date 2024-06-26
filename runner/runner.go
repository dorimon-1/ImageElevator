package runner

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type RunnerBase struct {
	ctx            context.Context
	sampleRate     time.Duration
	timer          *time.Timer
	runUploadChan  chan any
	resetTimerChan chan any
}

func NewRunnerBase(sampleRate time.Duration) RunnerBase {
	return RunnerBase{
		ctx:            context.Background(),
		sampleRate:     sampleRate,
		timer:          time.NewTimer(sampleRate),
		runUploadChan:  make(chan any, 1),
		resetTimerChan: make(chan any),
	}
}

type Runner interface {
	runnerBase() *RunnerBase
	Stop() error
	uploadImages() (int, error)
	pullFiles() ([]string, error)
}

func TriggerUpload(r Runner) error {
	select {
	case r.runnerBase().runUploadChan <- nil:
		return nil
	default:
		return errors.New("too many requests")
	}
}

func Start(r Runner) {
	go timerRoutine(r)
	go uploaderRoutine(r)
}

func uploaderRoutine(r Runner) {
	log.Debug().Msg("Image uploader routine started")
	for {
		select {
		case <-r.runnerBase().runUploadChan:
			imagesCount, err := r.uploadImages()
			if err != nil {
				log.Error().Msgf("Failed to upload images from Image Uploader => %s", err)
			}

			if imagesCount > 0 {
				log.Info().Msgf("Uploaded %d images", imagesCount)
			} else {
				log.Debug().Msg("No images uploaded")
			}

			r.runnerBase().resetTimerChan <- nil

		case <-r.runnerBase().ctx.Done():
			log.Debug().Msg("Closing Image Uploader")
			if err := r.Stop(); err != nil {
				log.Warn().Msgf("received an error while stopping runner => %s", err)
			}
			close(r.runnerBase().runUploadChan)
			close(r.runnerBase().resetTimerChan)
			return
		}
	}
}
