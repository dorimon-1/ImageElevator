package runner

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/Kjone1/imageElevator/ftp"
	"github.com/rs/zerolog/log"
)

type RunnerBase struct {
	ctx            context.Context
	sampleRate     time.Duration
	ftpClient      ftp.FTPClient
	timer          *time.Timer
	runUploadChan  chan any
	resetTimerChan chan any
	workingPath    string
	filePattern    string
	uploadedFiles  map[string]bool
}

func NewRunnerBase(sampleRate time.Duration, ftpClient ftp.FTPClient, workingPath, filePattern string, uploadedFiles map[string]bool) RunnerBase {
	return RunnerBase{
		ctx:            context.Background(),
		sampleRate:     sampleRate,
		ftpClient:      ftpClient,
		timer:          time.NewTimer(sampleRate),
		runUploadChan:  make(chan any, 1),
		resetTimerChan: make(chan any),
		workingPath:    workingPath,
		filePattern:    filePattern,
		uploadedFiles:  make(map[string]bool),
	}
}

type Runner interface {
	runnerBase() *RunnerBase
	Stop() error
	uploadImages() (int, error)
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

func loadCache(fileName string) map[string]bool {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Warn().Msgf("Couldn't find cache file %s creating one.", fileName)
		return make(map[string]bool)
	}
	var files map[string]bool
	if err := json.Unmarshal(data, &files); err != nil {
		log.Error().Msgf("Error reading file %s: %s", fileName, err)
		return make(map[string]bool)
	}

	return files
}

func saveCache(fileName string, files map[string]bool) error {
	data, err := json.Marshal(files)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
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
