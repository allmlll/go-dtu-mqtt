package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
	"strconv"
)

type cemsController struct {
}

var Cems cemsController

func (*cemsController) Get(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if !page.Check() {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", nil))
		return
	}

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
		c.JSON(http.StatusOK, util.ErrorMess("缺少code", err.Error()))
		return
	}

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	datas, err := service.Cems.Get(&page, interval, start, end, code)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("cems数据获取失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("cems数据获取成功", datas))
	}
}

func (*cemsController) BigScreen(c *gin.Context) {

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
		c.JSON(http.StatusOK, util.ErrorMess("缺少code", err.Error()))
		return
	}

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	datas, err := service.Cems.BigScreen(interval, start, end, code)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("cems数据获取失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("cems数据获取成功", datas))
	}
}

func (*cemsController) EveryCemsLatestData(c *gin.Context) {
	data, err := service.Cems.EveryCemsLatestData()
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("cems数据获取失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("cems数据获取成功", data))
	}
}
