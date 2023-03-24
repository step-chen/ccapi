package main

// CGO_ENABLED=0 go build -ldflags "-s -w" .
// CGO_ENABLED=0 GOOS=linux go build -o ./apps/ccapi -ldflags '-s -w --extldflags "-static -fpic"' main.go

import (
	"ccapi/internal/pkg/bookinfo"

	"github.com/gin-gonic/gin"
)

// curl http://localhost:8080/isbn/9787542679307
func main() {
	//gin.ForceConsoleColor()

	//f, _ := os.Create("ccapi.log")
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	router.GET("/isbn/:isbn", bookinfo.GetByIsbnFromDouban)

	router.Run(":8080")
}
