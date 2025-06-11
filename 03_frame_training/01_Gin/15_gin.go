package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os"
)

// 发送https请求
var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

func main() {
	router := gin.Default()
	os.Chdir("D:\\golang projects\\go_training\\03_frame_training\\01_Gin")
	router.Static("/assets", "./assets")
	router.SetHTMLTemplate(html)

	router.GET("/", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			// 这里是使用https主动推送静态资源到客户端
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(200, "https", gin.H{
			"status": "success",
		})
	})

	// Listen and Serve in https://127.0.0.1:8080
	router.RunTLS(":8080", "./keys/server.pem", "./keys/server.key")
}
