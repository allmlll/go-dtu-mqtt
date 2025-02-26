package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
	"ruitong-new-service/middleware"
)

func GetEngine() *gin.Engine {
	engine := gin.Default()

	engine.Use(middleware.Cors())                  //跨域
	engine.Use(gzip.Gzip(gzip.DefaultCompression)) //数据压缩
	engine.POST("/login", controller.User.Login)
	engine.GET("/project/get", controller.Project.Get)
	engine.GET("/file/checkApk", controller.File.CheckApk)
	engine.Use(middleware.JwtAuth(), middleware.ApiAuth()) //权限
	engine.GET("/latestData/bigScreenCems/ws", controller.LatestData.BigScreenCems)
	Device(engine)
	User(engine)
	Role(engine)
	Api(engine)
	Dtu(engine)
	Cems(engine)
	Log(engine)
	File(engine)
	BigScreen(engine)
	LatestData(engine)
	Group(engine)
	Center(engine)
	Project(engine)
	Menu(engine)
	Report(engine)
	Param(engine)
	WebSocket(engine)
	return engine
}
