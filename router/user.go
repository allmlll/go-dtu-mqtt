package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func User(engine *gin.Engine) {
	user := engine.Group("user")
	{
		user.POST("/register", controller.User.Register)
		user.PUT("/update", controller.User.Update)
		user.GET("/get", controller.User.Get)
		user.DELETE("/delete", controller.User.Delete)
		user.POST("/resetPassword", controller.User.ResetPass)
		user.POST("/changePassword", controller.User.ChangePass)
	}
}
