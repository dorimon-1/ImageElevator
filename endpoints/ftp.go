package endpoints

import (
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
	h.runner.TriggerUpload()
	c.String(200, "sync requested succesfully")
}
