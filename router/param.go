package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Param(engine *gin.Engine) {
	device := engine.Group("/param")
	{
		device.POST("/add", controller.Param.Add)
		device.DELETE("/delete", controller.Param.Delete)
		device.PUT("/update", controller.Param.Update)
		device.GET("/pageGet", controller.Param.Find)
	}
}
