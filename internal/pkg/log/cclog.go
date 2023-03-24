package cclog

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Log(format string, values ...any) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	f := fmt.Sprintf("[CCAPI] %s %s\n", now, format)
	fmt.Fprintf(gin.DefaultWriter, f, values...)
}

func LogErr(format string, values ...any) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	f := fmt.Sprintf("[CCAPI-ERROR] %s | %s\n", now, format)
	fmt.Fprintf(gin.DefaultWriter, f, values...)
}
