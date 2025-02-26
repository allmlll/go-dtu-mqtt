package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Api(engine *gin.Engine) {
	api := engine.Group("api")
	{
		api.GET("/pageGet", controller.Api.PaginationFind)
		api.GET("/get", controller.Api.Get)
	}
}
