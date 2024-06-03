package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func main() {
  server := gin.Default()
  server.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "am i..? am i realy alive..?",
    })
  })
  server.Run()
}
