package handler

import (
	"net/http"

	"github.com/Kjone1/imageElevator/elevator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	elevator elevator.Elevator
}

func NewHandler(elevator elevator.Elevator) *Handler {
	return &Handler{
		elevator: elevator,
	}
}

func (h *Handler) Sync(c *gin.Context) {
	if err := elevator.TriggerUpload(h.elevator); err != nil {
		c.String(http.StatusTooManyRequests, err.Error())
		return
	}

	c.String(http.StatusOK, "sync requested succesfully")
}

func (h *Handler) Health(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
