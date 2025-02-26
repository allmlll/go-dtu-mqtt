package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Device(engine *gin.Engine) {
	device := engine.Group("/device")
	{
		device.POST("/add", controller.Device.Add)
		device.DELETE("/delete", controller.Device.Delete)
		device.PUT("/update", controller.Device.Update)
		device.GET("/pageGet", controller.Device.PaginationFind)
	}
}
