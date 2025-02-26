package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
	"strconv"
)

type dtuController struct {
}

var Dtu dtuController

func (d *dtuController) Get(c *gin.Context) {
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

	datas, err := service.Dtu.Get(&page, interval, start, end, code)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取"+code+"dtu设备数据错误", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取"+code+"dtu设备数据成功", datas))
	}
}

func (d *dtuController) BigScreen(c *gin.Context) {
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

	datas, err := service.Dtu.BigScreen(interval, start, end, code)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取"+code+"dtu设备数据错误", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取"+code+"dtu设备数据成功", datas))
	}
}
