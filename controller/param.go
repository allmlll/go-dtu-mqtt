package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ruitong-new-service/model"
	"ruitong-new-service/service"
	"ruitong-new-service/util"
)

type paramController struct{}

var Param paramController

func (p *paramController) Add(c *gin.Context) {
	var param model.Param
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.Param.Add(param); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("添加参数失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("添加参数成功", data))
	}
}

func (p *paramController) Find(c *gin.Context) {
	var page util.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if !page.Check() {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", nil))
		return
	}

	if data, err := service.Param.Find(&page); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询参数失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询参数成功", data))
	}
}

func (p *paramController) Update(c *gin.Context) {
	var param model.Param
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.Param.Update(param); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("修改参数失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("修改参数成功", data))
	}
}

func (p *paramController) Delete(c *gin.Context) {
	var param model.Param
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.Param.Delete(param); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除参数失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("删除参数成功", data))
	}
}
