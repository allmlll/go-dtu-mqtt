package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type latestDataController struct {
}

var LatestData latestDataController

func (l latestDataController) Get(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
	}
}

func (l latestDataController) WebHome(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil || !page.Check() {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
	}

	if data, err := service.LatestData.WebHome(page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("错误", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("成功", data))
	}
}

func (l latestDataController) BgiScreenMap(c *gin.Context) {
	if data, err := service.LatestData.BgiScreenMap(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("错误", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("成功", data))
	}
}

func (l latestDataController) BigScreenCems(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		// 在这里可以添加更多的配置选项，比如检查Origin是否允许等
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("连接失败", err.Error()))
		return
	}

	// 读取客户端的消息
	_ = service.LatestData.BigScreenCems(conn)

	_ = conn.Close()
}
