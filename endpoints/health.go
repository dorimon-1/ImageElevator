package endpoints

import "github.com/gin-gonic/gin"

func Health(c *gin.Context) {
	c.String(200, "pong")
}
