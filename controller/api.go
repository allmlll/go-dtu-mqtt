package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type apiController struct {
}

var Api apiController

func (*apiController) PaginationFind(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if !page.Check() {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", nil))
		return
	}

	url := c.Query("url")
	method := c.Query("method")

	if data, err := service.Api.PaginationFind(page, url, method); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取api信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取api信息成功", data))
	}
}

func (*apiController) Get(c *gin.Context) {

	if data, err := service.Api.Get(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取api信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取api信息成功", data))
	}
}
