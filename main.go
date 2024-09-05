package main

// CGO_ENABLED=0 go build -ldflags "-s -w" .
// CGO_ENABLED=0 GOOS=linux go build -o ./apps/ccapi -ldflags '-s -w --extldflags "-static -fpic"' main.go

// curl http://localhost:5003/book/douban/isbn/9787542679307
// curl http://localhost:5003/book/douban/name/要命还是要灵魂/2

import (
	"github.com/step-chen/ccapi/internal/pkg/bookinfo"

	"github.com/gin-gonic/gin"
)

func main() {
	//gin.ForceConsoleColor()

	//f, _ := os.Create("ccapi.log")
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	router.GET("book/douban/isbn/:isbn", bookinfo.GetByIsbnFromDouban)
	router.GET("book/douban/name/:name/:count", bookinfo.GetByNameFromDouban)
	router.GET("book/douban/name/:name", bookinfo.GetByNameFromDouban)

	router.Run(":5003")
}
