package runner

import (
	"errors"
	"github.com/rs/zerolog/log"
)

type Runner interface {
	uploadImages() (int, error)
	pullFiles() ([]string, error)
	Start()
	TriggerUpload() error
}

func uploaderRoutine(r *DockerRunner) {
	log.Debug().Msg("Image uploader routine started")
	for {
		select {
		case <-r.runUploadChan:
			imagesCount, err := r.uploadImages()
			if err != nil {
				//FIXME: This is stupid!!! make it not an error somehow
				if err.Error() == "no new images were found" {
					log.Debug().Msg(err.Error())
				} else {
					log.Err(err).Msgf("Failed to upload images from Image Uploader")
				}
			}

			if imagesCount > 0 {
				log.Info().Msgf("Uploaded %d images", imagesCount)
			} else {
				log.Debug().Msg("No images uploaded")
			}

			r.resetTimerChan <- nil

		case <-r.ctx.Done():
			log.Debug().Msg("Closing Image Uploader")
			if err := r.ftpClient.Close(); err != nil {
				log.Err(err).Msg("Failed to close connection")
			}

			close(r.runUploadChan)
			close(r.resetTimerChan)
			return
		}
	}
}

func (r *DockerRunner) TriggerUpload() error {
	select {
	case r.runUploadChan <- nil:
		return nil
	default:
		return errors.New("too many requests")

	}
}
