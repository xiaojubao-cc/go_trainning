package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//静态文件资源
/*
	方法            适用场景            参数特性            目录列表控制          典型用途
  Static()        映射整个目录         路径字符串             自动禁用      托管 CSS/JS/图片等静态资源目录
 StaticFS()       自定义文件系统    http.FileSystem 接口     可自定义      嵌入式资源/内存文件系统/权限控制
StaticFile()      映射单个文件         文件路径字符串          不涉及       托管 favicon.ico 等独立文件
维度             		Static()          		StaticFS()              StaticFile()
路径匹配方式   前缀匹配（/assets/*filepath）          前缀匹配                  精确匹配
性能                 中等（需路径解析）                中等                 最高（直接映射）
安全性                自动防目录遍历           依赖 FileSystem 实现            无风险
扩展性               固定本地文件系统           支持任意文件系统               仅限本地文件
*/
func main() {
	router := gin.Default()
	//访问示例： /static/css/style.css → ./dist/css/style.css
	router.Static("/static", "./assets")
	// 对频繁访问的静态文件启用内存缓存

	router.StaticFS("/more_static", http.Dir("my_file_system"))

	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// 公共资源目录
	router.Static("/public", "./public")

	// 敏感文件单独控制
	router.StaticFile("/.well-known/acme-challenge/token", "./ssl/token.txt")

	/*curl -X POST http://localhost:8080/upload -F "file=@C:\Users\Administrator\Desktop\example.png" -F "name=wang" -F "email=1171264943@qq.com"*/
}
