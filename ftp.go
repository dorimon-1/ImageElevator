package main

import (
    "fmt"
    "time"
    "os"

    "github.com/secsy/goftp"
    "github.com/gin-gonic/gin"
)

const (
    ftpServerURL = "localhost"
    ftpServerPath = "/home/kj"
)

func ftpListEndpoint(c *gin.Context) {
    config := goftp.Config {
      User:               "kj",
      Password:           "1",
      ConnectionsPerHost: 10,
      Timeout:            10 * time.Second,
      Logger:             os.Stdout,
    }
    client, err := goftp.DialConfig(config,ftpServerURL)
    if err != nil {
        panic(err)
    }
    files, err := client.ReadDir(ftpServerPath)
    if err != nil {
        panic(err)
    }
    for _, file := range files {
        fmt.Println(file.Name())
    }
}
