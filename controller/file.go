package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/global"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type fileController struct {
}

var File fileController

func (*fileController) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		global.Log.Error(err.Error())
		c.JSON(http.StatusOK, util.ErrorMess("上传失败", err.Error()))
		return
	}
	url, err := service.File.Upload(file)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("上传失败", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("上传成功", url))
}

func (*fileController) CheckApk(c *gin.Context) {
	url, err := service.File.CheckApk()
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("最新版本检查失败", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.SuccessMess("最新版本检查成功", url))
}
