package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 设置请求头
var (
	expectedHost = "localhost:8080"
)

func SetSecurityHeaders(c *gin.Context) {
	//Host头验证
	if c.Request.Host != expectedHost {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid host header",
		})
	}
	//禁止页面被iframe嵌套
	c.Header("X-Frame-Options", "DENY")
	//内容安全策略防止XSS攻击
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	//启用浏览器XSS过滤
	c.Header("X-XSS-Protection", "1; mode=block")
	//强制使用HTTPS
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	//禁止Referrer信息泄露
	c.Header("Referrer-Policy", "strict-origin")
	//禁止MIME类型检测
	c.Header("X-Content-Type-Options", "nosniff")
	//禁用敏感设备功能
	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	//请求
	c.Next()
}
func main() {
	router := gin.Default()
	router.Use(SetSecurityHeaders)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello http header!",
		})
	})
	router.Run(":8080")
}
