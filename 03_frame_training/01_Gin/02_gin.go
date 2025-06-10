package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// 绑定前端复选框
type colorForm struct {
	Colors []string `form:"colors[]"`
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}
func boundSlice(c *gin.Context) {
	var form colorForm
	if err := c.ShouldBind(&form); err == nil {
		/*gin.H快速构建json响应*/
		c.JSON(http.StatusOK, gin.H{
			"colors": form.Colors,
		})
	}
}
func main() {
	engine := gin.Default()
	/*设置工作目录为项目根目录*/
	os.Chdir("D:\\golang projects\\go_training\\03_frame_training\\01_Gin")
	engine.LoadHTMLGlob("views/*")
	engine.GET("/", index)
	engine.POST("/", boundSlice)
	engine.Run(":8080")
}
