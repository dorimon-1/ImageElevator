package handler

import (
	"net/http"

	"github.com/Kjone1/imageElevator/runner"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	runner *runner.Runner
}

func NewHandler(runner *runner.Runner) *Handler {
	return &Handler{
		runner: runner,
	}
}

func (h *Handler) Sync(c *gin.Context) {
	if err := h.runner.TriggerUpload(); err != nil {
		c.String(http.StatusTooManyRequests, err.Error())
		return
	}

	c.String(http.StatusOK, "sync requested succesfully")
}

func (h *Handler) Health(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
