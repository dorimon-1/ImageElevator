package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/endpoints"
	"github.com/Kjone1/imageElevator/runner"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	if os.Getenv("GIN_MODE") != "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
	config.LoadConfig()
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGPIPE,
	)
	defer cancel()
	server := gin.Default()

	v1 := server.Group("/v1")
	runner := runner.NewRunner(ctx)
	handler := endpoints.NewHandler(runner)

	v1.GET("/ping", endpoints.Health)
	v1.GET("/sync", handler.Sync)

	if err := server.Run(); err != nil {
		log.Fatal().Msgf("failed to start server: %s", err)
	}
}
