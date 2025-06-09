package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Person struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

func boundStruct(c *gin.Context) {
	var person Person
	err := c.ShouldBind(&person)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.String(http.StatusOK, `boundStruct`)
}
func main() {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello gin",
		})
	})
	engine.GET("/boundStruct", boundStruct)
	engine.Run(":8080")
}
