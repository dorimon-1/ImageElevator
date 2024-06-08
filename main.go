package main

import (
	"log"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/endpoints"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadConfig()
}

func main() {
	server := gin.Default()

	v1 := server.Group("/v1")

	v1.GET("/ping", endpoints.Health)
	v1.GET("/sync", endpoints.FtpSync)

	if err := server.Run(); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
