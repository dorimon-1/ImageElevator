package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/endpoints"
	"github.com/gin-gonic/gin"
)

func init() {
	if os.Getenv("GIN_MODE") != "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
	config.LoadConfig()
}

func main() {
	server := gin.Default()

	v1 := server.Group("/v1")

	v1.GET("/ping", endpoints.Health)
	v1.GET("/sync", endpoints.FtpSync)

	if err := server.Run(); err != nil {
		log.Fatal().Msgf("failed to start server: %s", err)
	}
}
