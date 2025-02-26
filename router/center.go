package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Center(engine *gin.Engine) {
	center := engine.Group("/center")
	{
		center.POST("/add", controller.Center.Add)
		center.DELETE("/delete", controller.Center.Delete)
		center.PUT("/update", controller.Center.Update)
		center.GET("/get", controller.Center.Get)
		center.GET("/pageGet", controller.Center.PaginationFind)
	}
}
