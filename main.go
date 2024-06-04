package main

import (
	"github.com/gin-gonic/gin"
)

func init() {
	LoadConfig()
}

func main() {
	server := gin.Default()
	v1 := server.Group("/v1")
	{
		v1.GET("/ping", healthEndpoint)
		v1.GET("/list", ftpListEndpoint)
	}

	server.Run()
}
