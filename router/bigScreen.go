package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func BigScreen(engine *gin.Engine) {
	bigScreen := engine.Group("/bigScreen")
	{
		bigScreen.POST("/get", controller.BigScreen.Get)
		bigScreen.GET("/getAll", controller.BigScreen.GetAll)
		bigScreen.GET("/getCount", controller.BigScreen.GetCount)
		bigScreen.GET("/getHisStove", controller.BigScreen.GetHisStove)
		bigScreen.GET("/getReportData", controller.BigScreen.GetReportData)
		bigScreen.GET("/getProfilingIngredientsData", controller.BigScreen.GetProfilingIngredientsData)
	}
}
