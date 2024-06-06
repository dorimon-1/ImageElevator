package endpoints

import (
	"log"

	"github.com/Kjone1/imageElevator/ftp"
	"github.com/gin-gonic/gin"
)

func FtpSync(c *gin.Context) {
	client := ftp.Client()
	images, err := ftp.List(client)
	if err != nil {
		log.Printf("Reading FTP directory failed with error => %s", err)
		return
	}
	if images == nil {
		log.Printf("No new images were found")
		return
	}
	ftp.Pull(client, images)
}
