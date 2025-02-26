package router

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/controller"
)

func WebSocket(engine *gin.Engine) {
	webSocket := engine.Group("/ws")
	{
		webSocket.GET("/getGraphiteWSData", controller.BigScreen.GetGraphiteWSData)
		webSocket.GET("/getTunnelKilnWSData", controller.BigScreen.GetTunnelKilnWSData)
		webSocket.GET("/getLargePressesWSData", controller.BigScreen.GetLargePressesWSData)
		webSocket.GET("/getDippingWSData", controller.BigScreen.GetDippingWSData)
	}
}
