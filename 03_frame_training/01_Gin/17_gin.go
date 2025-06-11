package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//绑定

type LoginForm struct {
	User     string ` binding:"required" json:"username"`
	Password string ` binding:"required" json:"password"`
}

func main() {
	router := gin.Default()
	router.POST("/login", func(c *gin.Context) {
		var loginForm LoginForm
		/*form 标签对应使用shouldBindWith只能处理form-data/x-www-form-urlencoded*/
		//if err := c.ShouldBindWith(&loginForm, binding.Form); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": err.Error(),
		//	})
		//}
		/*
				ShouldBindBodyWithJSON与ShouldBindWithJSON区别：
			    前者适用于多次绑定,后者适用于一次绑定
		*/
		//请求示例:curl -v -X POST ^
		//-H "Content-Type: application/json" ^
		//-d "{\"username\":\"admin\",\"password\":\"admin\"}" ^
		//http://localhost:8080/login
		if err := c.ShouldBindBodyWithJSON(&loginForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		if loginForm.User == "admin" && loginForm.Password == "admin" {
			c.JSON(http.StatusOK, gin.H{
				"status": "login success",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "login failed",
			})
		}
	})
	router.Run(":8080")
}
