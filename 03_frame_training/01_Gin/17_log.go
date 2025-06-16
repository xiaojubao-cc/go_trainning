package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// 记录日志
func main() {
	logFile, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(logFile)
	router := gin.Default()
	router.GET("/recordLog", func(context *gin.Context) {
		context.String(200, "success")
	})
	router.Run(":8080")
}
