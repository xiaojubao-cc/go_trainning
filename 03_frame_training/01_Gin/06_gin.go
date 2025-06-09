package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 自定义日志格式
func customLogFormatter(param gin.LogFormatterParams) string {
	// 自定义日志输出格式
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

func main() {
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(customLogFormatter))
	//捕获应用运行时的 panic 并恢复服务，防止未处理的崩溃导致整个进程退出
	router.Use(gin.Recovery())
	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "hello customLogFormatter",
		})
	})
	router.Run(":8080")
}
