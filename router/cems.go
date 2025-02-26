package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Cems(engine *gin.Engine) {
	cems := engine.Group("/cems")
	{
		cems.GET("/get", controller.Cems.Get)
		cems.GET("/bigScreen", controller.Cems.BigScreen)
		cems.GET("/everyCemsLatestData", controller.Cems.EveryCemsLatestData)
	}
}
