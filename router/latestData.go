package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func LatestData(engine *gin.Engine) {
	group := engine.Group("/latestData")
	{
		group.GET("/webHome", controller.LatestData.WebHome)
		group.GET("/bigScreen/map", controller.LatestData.BgiScreenMap)
	}
}
