package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Project(engine *gin.Engine) {
	project := engine.Group("project")
	{

		project.PUT("/update", controller.Project.Update)
		project.POST("/add", controller.Project.Add)
		project.DELETE("/delete", controller.Project.Delete)
	}
}
