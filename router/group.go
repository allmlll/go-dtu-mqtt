package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Group(engine *gin.Engine) {
	group := engine.Group("/group")
	{
		group.POST("/add", controller.Group.Add)
		group.DELETE("/delete", controller.Group.Delete)
		group.PUT("/update", controller.Group.Update)
		group.GET("/get", controller.Group.Get)
		group.GET("/getCascade", controller.Group.GetCascade)
		group.GET("/pageGet", controller.Group.PaginationFind)
	}
}
