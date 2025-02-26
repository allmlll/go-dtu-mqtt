package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/model"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type roleController struct {
}

var Role roleController

func (*roleController) Add(c *gin.Context) {
	var role model.Role
	if err := c.Bind(&role); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Role.Add(&role)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("添加角色失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("添加角色成功", data))
}

func (*roleController) PaginationFind(c *gin.Context) {
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

	if data, err := service.Role.PaginationFind(&page, name); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询角色失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询角色成功", data))
	}
}

func (*roleController) Get(c *gin.Context) {

	if data, err := service.Role.Get(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询角色失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询角色成功", data))
	}
}

func (*roleController) Update(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.Role.Update(&role); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("更新角色失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("更新角色成功", data))
	}
}
func (*roleController) GetCascade(c *gin.Context) {
	if data, err := service.Role.GetCascade(); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("级联菜单查询失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("级联菜单查询成功", data))
	}
}

//func (roleController) Delete(c *gin.Context) {
//	var device model.Device
//	if err := c.ShouldBindJSON(&device); err != nil {
//		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
//		return
//	}
//
//	data, err := service.Device.Delete(&device)
//	if err != nil {
//		c.JSON(http.StatusOK, util.ErrorMess("删除失败", err.Error()))
//	}
//
//	err = service.LatestData.Delete(device.Code)
//	if err != nil {
//		c.JSON(http.StatusOK, util.ErrorMess("删除最新数据失败", err.Error()))
//	}
//
//	c.JSON(http.StatusOK, util.SuccessMess("删除成功", data))
//}
