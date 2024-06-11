package endpoints

import (
	"github.com/rs/zerolog/log"

	"github.com/Kjone1/imageElevator/ftp"
	"github.com/gin-gonic/gin"
)

func FtpSync(c *gin.Context) {
	client, err := ftp.Client()
	if err != nil {
		log.Error().Msgf("Unable to create FTP client with error => %s", err)
		return
	}
	images, err := ftp.List(client)
	if err != nil {
		log.Error().Msgf("Reading FTP directory failed with error => %s", err)
		return
	}
	if images == nil {
		log.Info().Msg("No new images were found")
		return
	}

	_ , err = ftp.Pull(client, images)
	if err != nil {
		log.Error().Msgf("Pulling images from FTP server failed with error => %s", err)
	}
}
