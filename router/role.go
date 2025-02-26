package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Role(engine *gin.Engine) {
	role := engine.Group("role")
	{
		role.POST("/add", controller.Role.Add)
		role.PUT("/update", controller.Role.Update)
		role.GET("/get", controller.Role.Get)
		role.GET("/pageGet", controller.Role.PaginationFind)
		role.GET("/getCascade", controller.Role.GetCascade)
	}
}
