package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func Menu(engine *gin.Engine) {
	menu := engine.Group("menu")
	{
		menu.GET("/getMenu", controller.MenuController.GetMenu)
	}
}
