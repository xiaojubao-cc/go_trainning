package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// 在handler.Func中使用携程需要针对上下文进行拷贝
func main() {

	router := gin.Default()
	router.GET("/varify", func(context *gin.Context) {
		copy := context.Copy()
		/*这里携程可以打印出来因为主携程未中断*/
		go func() {
			time.Sleep(10 * time.Second)
			fmt.Printf("delay 10 second get url：%s", copy.Request.URL.String())
		}()
		context.JSON(200, gin.H{
			"message": "hello varify",
		})
	})
	router.Run(":8080")
}
