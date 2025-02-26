package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Dtu(engine *gin.Engine) {
	dtu := engine.Group("/dtu")
	{
		dtu.GET("/get", controller.Dtu.Get)
		dtu.GET("/bigScreen", controller.Dtu.BigScreen)
	}
}
