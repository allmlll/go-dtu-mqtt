package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/model"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type projectController struct {
}

var Project projectController

func (*projectController) Get(c *gin.Context) {
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
	if data, err := service.Project.Get(&page, url); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取项目信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取项目信息成功", data))
	}
}

func (*projectController) Update(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.Project.Update(&project); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("更新项目信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("更新项目信息成功", data))
	}
}

func (*projectController) Add(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.Project.Add(&project); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("添加项目信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("添加项目信息成功", data))
	}
}

func (*projectController) Delete(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.Project.Delete(&project); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除项目信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("删除项目信息成功", data))
	}
}
