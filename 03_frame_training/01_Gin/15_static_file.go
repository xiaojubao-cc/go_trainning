package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
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
//	// 公共资源目录
//	router.Static("/public", "./public")
//	// 敏感文件单独控制
//	router.StaticFile("/.well-known/acme-challenge/token", "./ssl/token.txt")
//使用 var f embed.FS 的核心优势是 将资源文件编译进程序，实现零依赖部署和统一管理，非常适合现代云原生应用、CLI 工具、小型 Web 应用等场景
/*将assets/* templates/*路径下的资源打包进代码*/
//go:embed assets/* templates/*
var fs embed.FS

func main() {
	router := gin.Default()
	os.Chdir("D:\\golang projects\\go_training\\03_frame_training\\01_Gin")
	tl := template.Must(template.New("").ParseFS(fs, "templates/*.tmpl", "templates/foo/*.tmpl"))
	router.SetHTMLTemplate(tl)
	/*
			assets 用于嵌入式资源(嵌入代码的资源)，通常是程序内部使用。
		    public 用于对外提供静态文件服务，客户端可直接访问
	*/
	//访问示例： /css/style.css → ./dist/css/style.css
	//router.Static("/static", "./assets")
	// 对频繁访问的静态文件启用内存缓存
	router.StaticFS("/public", http.FS(fs))

	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "static Resource",
		})
	})

	router.GET("/foo", func(context *gin.Context) {
		context.HTML(http.StatusOK, "bar.tmpl", gin.H{
			"title": "static Resource",
		})
	})

	router.Run(":8080")
}
