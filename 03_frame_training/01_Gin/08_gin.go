package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"time"
)

type Booking struct {
	CheckIn  time.Time `uri:"check_in" binding:"required,customValidator" time_format:"2006-01-02"`
	CheckOut time.Time `uri:"check_out" binding:"required,gtfield=CheckIn,customValidator" time_format:"2006-01-02"`
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

func getBooking(c *gin.Context) {
	var booking Booking
	if err := c.ShouldBindUri(&booking); err == nil {
		c.JSON(200, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(400, gin.H{"error": err.Error()})
	}
}
func main() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("customValidator", customValidator())
	}
	router.GET("/booking/:check_in/:check_out", getBooking)
	router.Run(":8080")
}
