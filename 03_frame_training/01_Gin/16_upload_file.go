package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

// 上传文件	/*curl -X POST http://localhost:8080/upload -F "file=@C:\Users\Administrator\Desktop\example.png" -F "name=wang" -F "email=1171264943@qq.com"*/
func main() {
	router := gin.Default()
	// 8 * 2^20  = 8MB
	router.MaxMultipartMemory = 8 << 20
	//设置目录为当前目录
	os.Chdir("D:\\golang projects\\go_training\\03_frame_training\\01_Gin")
	/*
			当用户在浏览器中访问根路径（如 http://localhost:8080/）时，
			Gin 框架会从 ./public 目录下查找并返回对应的静态文件（例如 index.html、图片、CSS 文件等）。
		    如果 ./public 目录下存在 index.html 文件，那么访问根路径时会自动显示该页面。
	*/
	router.Static("/", "./public")
	router.POST("/singleFileUpload", singleFileUpload)
	router.POST("/multiFileUpload", multiFileUpload)
	router.Run(":8080")
}

// 单个文件上传
func singleFileUpload(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileName := filepath.Base(file.Filename)
	if errs := context.SaveUploadedFile(file, fileName); errs != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": errs.Error(),
		})
		return
	}
	context.String(http.StatusOK, "File uploaded successfully")
}

// 多个文件上传
func multiFileUpload(context *gin.Context) {
	name := context.PostForm("name")
	email := context.PostForm("email")
	multipartForm, err := context.MultipartForm()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	files := multipartForm.File["files"]
	for _, file := range files {
		fileName := filepath.Base(file.Filename)
		if err := context.SaveUploadedFile(file, fileName); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	context.String(http.StatusOK, "Files uploaded successfully name:%s,email:%s", name, email)
}
