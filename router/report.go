package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Report(engine *gin.Engine) {
	report := engine.Group("/report")
	{
		report.GET("/get", controller.Report.Get)
	}
}
