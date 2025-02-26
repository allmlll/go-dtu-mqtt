package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"ruitong-new-service/service"
	"ruitong-new-service/util"
	"strconv"

	"net/http"
)

type bigScreenController struct {
}

var BigScreen bigScreenController

func (*bigScreenController) Get(c *gin.Context) {
	r := service.BigScreen.Res

	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	all, err := service.BigScreen.Get(r)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("获取成功", all))
}

func (*bigScreenController) GetAll(c *gin.Context) {
	getAll, err := service.BigScreen.GetAll()
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("获取成功", getAll))
}

func (*bigScreenController) GetCount(c *gin.Context) {
	data, err := service.BigScreen.GetCount()
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("获取成功", data))
}

func (*bigScreenController) GetGraphiteWSData(c *gin.Context) {
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
	_ = service.BigScreen.GetGraphiteWSData(conn)

	_ = conn.Close()
}

func (*bigScreenController) GetTunnelKilnWSData(c *gin.Context) {
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
	_ = service.BigScreen.GetTunnelKilnWSData(conn)

	_ = conn.Close()
}

func (*bigScreenController) GetLargePressesWSData(c *gin.Context) {
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
	_ = service.BigScreen.GetLargePressesWSData(conn)

	_ = conn.Close()
}

func (*bigScreenController) GetDippingWSData(c *gin.Context) {
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
	_ = service.BigScreen.GetDippingWSData(conn)

	_ = conn.Close()
}

func (*bigScreenController) GetHisStove(c *gin.Context) {
	data, err := service.BigScreen.GetHisStove()
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("获取成功", data))
}

func (*bigScreenController) GetReportData(c *gin.Context) {
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	interval, err := strconv.Atoi(c.DefaultQuery("interval", "0"))
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if !(interval >= 0 && interval <= 3) {
		c.JSON(http.StatusOK, util.ErrorMess("时间粒度错误", nil))
		return
	}

	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusOK, util.ErrorMess("缺少code", nil))
		return
	}
	data, err := service.BigScreen.GetReportData(start, end, code, interval)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("获取成功", data))
}

func (*bigScreenController) GetProfilingIngredientsData(c *gin.Context) {
	data, err := service.BigScreen.GetProfilingIngredientsData()
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("获取成功", data))
}
