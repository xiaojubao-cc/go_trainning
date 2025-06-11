package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// 绑定汇总
/*
curl命令总结：
参数                  说明
	-X GET        		指定 HTTP 方法为 GET（可省略，默认 GET）
	-H            		设置请求头（如 -H "Content-Type: application/json"）
	-d            		发送 POST 请求体数据（适用于 POST/PUT 方法）(默认 application/x-www-form-urlencoded)
	-G            		强制将 -d 参数转换为 URL 查询参数（GET 方法专用）
	-v                  输出详细信息
	-o                  将响应输出到文件
	-O                  下载文件并保留原文件
	-k                  忽略SSL证书验证
	-u                  添加认证信息(格式：username:password)
	-b                  发送cookie
	-A                  设置User-Agent
	--compressed        压缩响应
	-T                  上传文件
	-L                  跟随重定向
	--form              上传表单/文件
	--data-urlencode    自动 URL 编码参数值（处理特殊字符时使用）
*/
//ShouldBindBodyWithJSON与ShouldBindWithJSON区别：	前者适用于同一次请求中多次绑定,后者适用于一次绑定
func main() {
	router := gin.Default()
	routerGroup := router.Group("/api")
	GetMethodBindingGroup(routerGroup)
	PostMethodBindingGroup(routerGroup)
	router.Run(":8080")
}

func GetMethodBindingGroup(rp *gin.RouterGroup) {
	routerGroup := rp.Group("/getMethod")
	/*http://localhost:8080/api/getMethod/ShouldBindQuery?name=appleboy&age=17*/
	routerGroup.GET("/ShouldBindQuery", shouldBindQuery)
	/*http://localhost:8080/api/getMethod/ShouldBindUri/appleboy/17*/
	routerGroup.GET("/ShouldBindUri/:name/:age", shouldBindUri)
	/*curl -X GET http://localhost:8080/api/getMethod/queryKey?name="wang?age=18"*/
	routerGroup.GET("/queryKey", noBindingQueryKey)
	/*http://localhost:8080/queryMap?user[name]=wg*/
	routerGroup.GET("/queryMap", noBindQueryMap)
	/*http://localhost:8080/api/getMethod/urlParam/wang/nihao/18这里*age获取到的是/nihao/18*/
	routerGroup.GET("/urlParam/:name/*age", urlParam)

}

/*获取路径参数*/
func urlParam(context *gin.Context) {
	name := context.Param("name")
	age := context.Param("age")
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("urlParam binding struct name:%s,age:%s", name, age),
	})
}

func PostMethodBindingGroup(rp *gin.RouterGroup) {
	routerGroup := rp.Group("/postMethod")
	/*curl -v -X POST -H "Content-Type: application/json" -d "{\"name\":\"admin\",\"age\":18}" http://localhost:8080/api/postMethod/ShouldBindWith*/
	routerGroup.POST("/ShouldBindWithJSON", shouldBindWithJSON)
	/*curl -v -X POST -H "x-www-form-urlencoded" -d "name=\"wang"\" -d "age=18" http://localhost:8080/api/postMethod/ShouldBindWithForm*/
	/*binding.Query调用示例:curl -G "http://localhost:8080/booking"--data-urlencode "check_in=2024-12-01"*/
	routerGroup.POST("/ShouldBindWithForm", shouldBindWithForm)
	/*curl -v -X POST -d user="li" http://localhost:8080/form*/
	routerGroup.POST("/postForm", noBindingFormKey)
	/*curl -v -X POST --data-raw "user[name]=li&user[age]=18" http://localhost:8080/formMap*/
	routerGroup.POST("/postFormMap", noBindingFormMap)
	/*curl -v -H "Content-Type: application/x-www-form-urlencoded" --data-raw "name=manu&message=this_is_great" "http://localhost:8080/api/postMethod/queryAndForm/?id=1234&page=1"*/
	routerGroup.POST("/queryAndForm", queryAndFormParam)
}

/*
	场景                                             选择方案
	快速获取 URL 或表单中的简单 map               QueryMap / PostFormMap
	处理复杂 JSON/XML 请求体                    ShouldBindWith + 结构体
	需要参数校验(如必填、格式)                    ShouldBindWith + binding 标签
	混合参数（如同时处理 Query 和 Body）          ShouldBind 系列（自动推断来源）
	格式x-www-form-urlencoded                 JSON/XML/YAML 等
*/
/*获取参数*/
func queryAndFormParam(context *gin.Context) {
	id := context.Query("id")
	page := context.Query("page")
	limit := context.DefaultQuery("limit", "10")
	username := context.PostForm("username")
	password := context.PostForm("password")
	number := context.DefaultPostForm("number", "12")
	sprintf := fmt.Sprintf("queryAndFormParam  id:%s,page:%s,limit:%s,username:%s,password:%s,number:%s", id, page, limit, username, password, number)
	context.String(http.StatusOK, sprintf)
}

/*绑定路径*/
func shouldBindUri(context *gin.Context) {
	var person struct {
		Name string `uri:"name" binding:"required"`
		Age  int    `uri:"age" binding:"required"`
	}
	/*该方法需要映射到struct并且需要有uri标签*/
	err := context.ShouldBindUri(&person)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ShouldBindUri binding struct name:%s,age:%d", person.Name, person.Age),
	})
}

/*绑定实体*/
func shouldBindQuery(context *gin.Context) {
	var person struct {
		Name string `form:"name" binding:"required"`
		Age  int    `form:"age" binding:"required"`
	}
	/*该方法需要映射到struct并且需要有form标签*/
	err := context.ShouldBindQuery(&person)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("shouldBindQuery binding struct name:%s,age:%d", person.Name, person.Age),
	})
}

/*绑定实体*/
/*
	binding.JSON 对应tag标签json,binding.XML对应的tag标签是xml
	binding.Uri对应的tag标签是uri,binding.Query(参数)/binding.Form(表单)对应的tag标签是form
*/
func shouldBindWithJSON(context *gin.Context) {
	var person struct {
		Name string `json:"name" binding:"required"`
		Age  int    `json:"age" binding:"required"`
	}
	err := context.ShouldBindWith(&person, binding.JSON)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("shouldBindWithJSON binding struct name:%s,age:%d", person.Name, person.Age),
	})
}

/*绑定form表单格式为x-www-form-urlencoded*/
func shouldBindWithForm(context *gin.Context) {
	var person struct {
		Name string `form:"name" binding:"required"`
		Age  int    `form:"age" binding:"required"`
	}
	err := context.ShouldBindWith(&person, binding.Form)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ShouldBindWithForm binding struct name:%s,age:%d", person.Name, person.Age),
	})
}

/*不需要绑定直接获取路径值Map不需要验证*/
func noBindQueryMap(context *gin.Context) {
	queryMap := context.QueryMap("user")
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("path variable is %s", queryMap),
	})
}

/*不需要绑定直接获取路径值不需要验证*/
func noBindingQueryKey(context *gin.Context) {
	name := context.Query("name")
	//设置默认值
	name = context.DefaultQuery("name", "default")
	age := context.Query("age")
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("queryKey nobinding struct name:%s,age:%s", name, age),
	})
}

/*不绑定实体*/
func noBindingFormKey(context *gin.Context) {
	value := context.PostForm("name")
	//未传时设置默认值
	value = context.DefaultPostForm("name", "default")
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("postForm nobinding struct name:%s", value),
	})
}
func noBindingFormMap(context *gin.Context) {
	postFormMap := context.PostFormMap("user")
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("path variable is %s", postFormMap),
	})
}
