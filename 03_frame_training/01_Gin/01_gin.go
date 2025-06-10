package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Person struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

func boundStruct(c *gin.Context) {
	var person Person
	/*绑定参数调用curl -X GET "localhost:8080/boundStruct?name=appleboy&age=17"*/
	err := c.ShouldBindQuery(&person)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	str := fmt.Sprintf("name: %s, age: %d", person.Name, person.Age)
	c.String(http.StatusOK, str)
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
