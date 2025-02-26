package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func File(engine *gin.Engine) {
	engine.POST("/image/upload", controller.File.Upload)
}
