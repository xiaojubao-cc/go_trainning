package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 绑定Uri

type UriPerson struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func boundUri(c *gin.Context) {
	var person UriPerson
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, person)
}
func main() {
	router := gin.Default()
	router.GET("/:name/:id", boundUri)
	router.Run(":8080")
}
