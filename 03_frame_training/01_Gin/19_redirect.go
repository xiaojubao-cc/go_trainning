package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// POST重定向状态码：http.StatusTemporaryRedirect:会保持原始请求方法（POST）进行重定向;StatusPermanentRedirect:永久重定向
func main() {
	//1.Get重定向
	router := gin.Default()
	router.GET("/getForward", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})
	//2.Post重定向curl -v -L -X POST http://localhost:8080/postForward -L跟随重定向
	router.POST("/postForward", func(context *gin.Context) {
		context.Redirect(http.StatusTemporaryRedirect, "/destUrl")
	})
	router.POST("/destUrl", func(context *gin.Context) {
		context.JSON(http.StatusOK, "post forward destUrl")
	})
	//3.router重定向curl -v  http://localhost:8080/routerForward保持源路径和目标路径都是一种方法GET/POST
	router.POST("/routerForward", func(context *gin.Context) {
		context.Request.URL.Path = "/getRouterDestUrl"
		router.HandleContext(context)
	})
	router.POST("/getRouterDestUrl", func(context *gin.Context) {
		context.JSON(http.StatusOK, "post forward RouterDestUrl")
	})
	router.Run(":8080")
}
