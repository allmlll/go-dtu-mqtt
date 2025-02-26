package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type logController struct {
}

var Log logController

func (*logController) PaginationFind(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if !page.Check() {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", nil))
		return
	}

	if data, err := service.Log.PaginationFind(page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取log信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取log信息成功", data))
	}
}

func (*logController) DeleteAll(c *gin.Context) {
	if data, err := service.Log.DeleteAll(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除log信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("删除log信息成功", data))
	}
}
