package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 自定义http请求
func main() {
	router := gin.Default()
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second, // 请求头读取超时
		WriteTimeout:   10 * time.Second, // 响应写入超时
		IdleTimeout:    30 * time.Second, // 空闲连接超时
		MaxHeaderBytes: 1 << 20,          // 最大请求头 1MB
	}
	server.ListenAndServe()
}
