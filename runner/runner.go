package runner

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type BasicRunner struct {
	ctx            context.Context
	sampleRate     time.Duration
	timer          *time.Timer
	runUploadChan  chan any
	resetTimerChan chan any
}

type Runner interface {
	basicRunner() *BasicRunner
	Stop() error
	uploadImages() (int, error)
	pullFiles() ([]string, error)
}

func TriggerUpload(r Runner) error {
	select {
	case r.basicRunner().runUploadChan <- nil:
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
		case <-r.basicRunner().runUploadChan:
			imagesCount, err := r.uploadImages()
			if err != nil {
				log.Error().Msgf("Failed to upload images from Image Uploader => %s", err)
			}

			if imagesCount > 0 {
				log.Info().Msgf("Uploaded %d images", imagesCount)
			} else {
				log.Debug().Msg("No images uploaded")
			}

			r.basicRunner().resetTimerChan <- nil

		case <-r.basicRunner().ctx.Done():
			log.Debug().Msg("Closing Image Uploader")
			close(r.basicRunner().runUploadChan)
			close(r.basicRunner().resetTimerChan)
			return
		}
	}
}
