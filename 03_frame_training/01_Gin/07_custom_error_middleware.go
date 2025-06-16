package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*自定义错误中间件*/

func customErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		//执行http请求调用
		context.Next()
		if len(context.Errors) > 0 {
			//获取最后的错误信息
			lastError := context.Errors.Last().Err
			context.JSON(http.StatusInternalServerError, gin.H{
				"code":    "1001",
				"message": lastError.Error(),
			})
		}
	}
}

func main() {
	router := gin.Default()
	//挂载
	router.Use(customErrorHandler())
	router.GET("/err", func(c *gin.Context) {
		//这里模拟一个错误
		c.Error(errors.New("this is a error"))
	})
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
