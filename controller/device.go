package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/model"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
	"time"
)

type deviceController struct {
}

var Device deviceController

func (*deviceController) Add(c *gin.Context) {
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if device.CheckTime == 0 {
		device.CheckTime = 5
	}
	if device.TimeInterval == nil || len(device.TimeInterval) == 0 {
		device.TimeInterval = []int{0}
	}
	if device.Sort == nil {
		device.Sort = []string{}
	}
	if device.ShowKeys == nil {
		device.ShowKeys = []string{}
	}
	data, err := service.Device.Add(&device)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("添加设备失败", err.Error()))
		return
	}
	// 最新数据表建立数据
	err = service.LatestData.Create(model.LatestData{
		Name:         device.Name,
		Code:         device.Code,
		Type:         device.Type,
		Data:         []model.Data{},
		TimeInterval: device.TimeInterval,
		CheckTime:    device.CheckTime,
		Status:       "在线",
		UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("添加最新数据失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("添加设备成功", data))
}

func (*deviceController) PaginationFind(c *gin.Context) {
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
	Type := c.Query("type")
	id := c.Query("id")

	if data, err := service.Device.PaginationFind(&page, name, Type, id); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询设备失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询设备成功", data))
	}
}

func (*deviceController) Update(c *gin.Context) {
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if device.CheckTime < 5 {
		device.CheckTime = 5
	}
	if device.TimeInterval == nil || len(device.TimeInterval) == 0 {
		device.TimeInterval = []int{0}
	}
	if device.Sort == nil {
		device.Sort = []string{}
	}
	if device.ShowKeys == nil {
		device.ShowKeys = []string{}
	}
	data, err := service.Device.Update(&device)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("更新设备失败", err.Error()))
		return
	}

	err = service.LatestData.Update(device.Code, device.Name, &device.TimeInterval, device.CheckTime)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("更新设备失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("更新设备成功", data))
}

func (*deviceController) Delete(c *gin.Context) {
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	data, err := service.Device.Delete(&device)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除设备失败", err.Error()))
		return
	}

	err = service.LatestData.Delete(device.Code)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除最新数据失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.SuccessMess("删除设备成功", data))
}
