package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type reportController struct {
}

var Report reportController

func (r *reportController) Get(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if !page.Check() {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", nil))
		return
	}

	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusOK, util.ErrorMess("缺少code", nil))
		return
	}

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	datas, err := service.Report.Get(&page, start, end, code)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取"+code+"设备报表数据错误", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取"+code+"设备报表数据成功", datas))
	}
}
