package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*自定义路由格式*/
func main() {
	router := gin.Default()
	/*这里相当于赋值给DebugPrintRouteFunc覆盖原本的DebugPrintRouteFunc函数,类似重写*/
	//gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	//	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	//}
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "status")
	})
	router.POST("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, "foo")
	})
	router.Run(":8080")
}
