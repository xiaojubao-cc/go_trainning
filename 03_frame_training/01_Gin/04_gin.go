package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// cookie
func cookieOperate(c *gin.Context) {
	cookie, err := c.Cookie("gin_cookie")
	if err != nil {
		/*cookie中各个参数释义:
		1. name: cookie的名称
		2. value: cookie的值
		3. maxAge: cookie的过期时间，单位是秒，默认是-1，表示关闭浏览器就过期
		4. path: cookie的路径，默认是"/"，表示根路径, "/api" 表示仅 /api 路径下生效。
		5. domain: cookie的域名，默认是空字符串，表示当前域名下生效。
		6. secure: cookie是否只允许https请求使用，默认是false
		7. httpOnly: true表示cookie只能通过http请求使用，不能通过js脚本访问,防XSS攻击。
		*/
		c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
	}
	fmt.Printf("current cookie is %s", cookie)
}
func main() {
	router := gin.Default()
	router.GET("/cookie", cookieOperate)
	router.Run(":8080")
}
