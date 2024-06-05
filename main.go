package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func healthEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
func main() {
	server := gin.Default()
	client, err := ftpConnect()
	if err != nil {
		log.Fatalln(err)
	}

	v1 := server.Group("/v1")
	{
		v1.GET("/ping", healthEndpoint)
		v1.GET("/list", client.ftpListEndpoint)
	}
	server.Run()
}
