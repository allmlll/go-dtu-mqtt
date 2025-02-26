package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Log(engine *gin.Engine) {
	log := engine.Group("log")
	{
		log.GET("/get", controller.Log.PaginationFind)
		log.DELETE("/deleteAll", controller.Log.DeleteAll)
	}
}
