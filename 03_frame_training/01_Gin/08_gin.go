package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"time"
)

type Booking struct {
	CheckIn  time.Time `form:"check_in" json:"check_in" binding:"required,customValidator" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" json:"check_out" binding:"required,gtfield=CheckIn,customValidator" time_format:"2006-01-02"`
	Name     string    `form:"name" json:"name" binding:"required,customLengthValidator"`
}

// 自定义验证器
func customValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if ok {
			today := time.Now()
			if today.After(date) {
				return false
			}
		}
		return true
	}
}
func customLengthValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		content, ok := fl.Field().Interface().(string)
		if ok && len(content) > 5 {
			return false
		}
		return true
	}
}

// 自定义结构体验证
func customStructValidator(sl validator.StructLevel) {
	if booking, ok := sl.Current().Interface().(Booking); ok {
		if booking.CheckIn.After(booking.CheckOut) {
			sl.ReportError(booking.CheckIn, "check_in", "CheckIn", "checkInAfterCheckOut", "")
		}
	}
}

func getBooking(c *gin.Context) {
	var booking Booking
	/*binding.JSON 对应tag标签json,binding.XML对应的tag标签是xml
	  binding.Uri对应的tag标签是uri,binding.Query(参数)/binding.Form(表单)对应的tag标签是form
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
	--form              上传表单/文件
	--data-urlencode    自动 URL 编码参数值（处理特殊字符时使用）
	  binding.Query调用示例:
	  curl -G "http://localhost:8080/booking" \
	  --data-urlencode "check_in=2024-12-01" \
	  --data-urlencode "check_out=2024-12-05" \
	  --data-urlencode "name=Lucy"
	  binding.Form调用示例:
	  curl -X POST http://localhost:8080/booking \
	  -d "check_in=2024-01-01" \
	  -d "check_out=2024-01-05" \
	  -d "name=John"
	*/
	/*shouldBindWith显式绑定,shouldBind自动检测,前者性能更高,精确控制,后者性能低更通用*/
	if err := c.ShouldBindWith(&booking, binding.Form); err == nil {
		c.JSON(200, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(400, gin.H{"error": err.Error()})
	}
}
func main() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("customValidator", customValidator())
		v.RegisterValidation("customLengthValidator", customLengthValidator())
		v.RegisterStructValidation(customStructValidator, Booking{})
	}
	router.POST("/booking", getBooking)
	router.Run(":8080")
}
