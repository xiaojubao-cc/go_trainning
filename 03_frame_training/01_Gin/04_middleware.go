package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

//自定义中间件

type CustomMiddleware = gin.HandlerFunc

func customMiddleware() CustomMiddleware {
	return func(context *gin.Context) {
		t := time.Now()
		// 设置 example 变量
		context.Set("example", "12345")
		// 请求前
		context.Next()
		// 请求后
		latency := time.Since(t)
		log.Print(latency)
		// 获取发送的 status
		status := context.Writer.Status()
		log.Println(status)
	}
}
func middlewareChain(router *gin.Engine, middlewares ...CustomMiddleware) {
	for _, middleware := range middlewares {
		router.Use(middleware)
	}
}
func main() {
	router := gin.New()
	//执行函数的调用
	middlewareChain(router, customMiddleware(), gin.Recovery())
	router.GET("/", func(context *gin.Context) {
		/*强制获取指定键值*/
		example := context.MustGet("example").(string)
		fmt.Printf("example: %s", example)
	})
	router.Run(":8080")
}
