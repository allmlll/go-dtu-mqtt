package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/model"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type groupController struct {
}

var Group groupController

func (*groupController) Add(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Group.Add(&group)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("添加分组失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("添加分组成功", data))
}

func (*groupController) PaginationFind(c *gin.Context) {
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
	id := c.Query("id")

	if data, err := service.Group.PaginationFind(&page, name, id); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("分页查询分组失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("分页查询分组成功", data))
	}
}

func (*groupController) Get(c *gin.Context) {
	if data, err := service.Group.Get(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询分组失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询分组成功", data))
	}
}
func (*groupController) GetCascade(c *gin.Context) {
	if data, err := service.Group.GetCascade(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询级联失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询级联成功", data))
	}
}

func (*groupController) Update(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Group.Update(&group)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("更新分组失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("更新分组成功", data))
}

func (*groupController) Delete(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Group.Delete(&group)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除分组失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("删除分组成功", data))
}
