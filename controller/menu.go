package controller

import (
	"github.com/gin-gonic/gin"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type menuController struct{}

var MenuController menuController

func (*menuController) GetMenu(c *gin.Context) {
	id, _ := c.Get("Id")
	if data, err := service.Menu.GetMenu(id); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("查询菜单失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("查询菜单成功", data))
	}
}
