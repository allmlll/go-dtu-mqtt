package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/model"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type centerController struct {
}

var Center centerController

func (*centerController) Add(c *gin.Context) {
	var center model.Center
	if err := c.ShouldBindJSON(&center); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Center.Add(&center)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("添加中心失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("添加中心成功", data))
}

func (*centerController) PaginationFind(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if !page.Check() {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", nil))
		return
	}
	name := c.Query("name")

	if data, err := service.Center.PaginationFind(&page, name); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("分页查询中心失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("分页查询中心成功", data))
	}
}

func (*centerController) Get(c *gin.Context) {
	if data, err := service.Center.Get(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询中心失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询中心成功", data))
	}
}

func (*centerController) Update(c *gin.Context) {
	var center model.Center
	if err := c.ShouldBindJSON(&center); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Center.Update(&center)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("更新中心失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("更新中心成功", data))
}

func (*centerController) Delete(c *gin.Context) {
	var center model.Center
	if err := c.ShouldBindJSON(&center); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Center.Delete(&center)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除中心失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("删除中心成功", data))
}
