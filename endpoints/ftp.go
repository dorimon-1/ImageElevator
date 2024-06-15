package endpoints

import (
	"github.com/rs/zerolog/log"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/ftp"
	"github.com/gin-gonic/gin"
)

// TODO: also rename file to sync.go
func Sync(c *gin.Context) {
	client, err := ftp.Client()
	if err != nil {
		log.Error().Msgf("Unable to create FTP client with error => %s", err)
		return
	}
	//TODO: make pattern an environement variable
	images, err := ftp.List(client, config.FtpConfig().FtpServerPath, "^int-.*-docker$")
	if err != nil {
		log.Error().Msgf("Reading FTP directory failed with error => %s", err)
		return
	}
	if images == nil {
		log.Info().Msg("No new images were found")
		return
	}

	_, err = ftp.Pull(client, images...)
	if err != nil {
		log.Error().Msgf("Pulling images from FTP server failed with error => %s", err)
	}
}
