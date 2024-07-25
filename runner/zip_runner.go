package runner

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/ftp"
	"github.com/rs/zerolog/log"
)

type ZipRunner struct {
	RunnerBase
	destinationPath string
}

func NewZipRunner(ctx context.Context, ftpClient ftp.FTPClient, runnerConfig *config.RunnerConfiguration, workingPath, filePattern, destPath string) *ZipRunner {
	return &ZipRunner{
		RunnerBase:      NewRunnerBase(runnerConfig.SampleRateInMinutes, ftpClient, workingPath, filePattern),
		destinationPath: destPath,
	}
}

func (r *ZipRunner) runnerBase() *RunnerBase {
	return &r.RunnerBase
}
func (r *ZipRunner) Stop() error {
	return nil
}
func (r *ZipRunner) uploadImages() (int, error) {
	count := 0
	zipFiles, err := r.pullFiles()
	if err != nil {
		return 0, err
	}

	for _, file := range zipFiles {
		fileName := path.Base(file)
		if fileName == "/" || fileName == "." {
			log.Error().Msg("Couldn't find file")
			continue
		}
		dest := r.destinationPath + "/" + fileName
		if err := copyFiles(file, dest); err != nil {
			log.Error().Err(err)
			continue
		}

		count++
	}

	return count, nil
}

// NOTE: destinationPath should be the path including the file name!!!
func copyFiles(file string, destinationPath string) error {
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	return nil
}
