package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ruitong-new-service/model"
	"ruitong-new-service/service"
	"ruitong-new-service/util"

	"net/http"
)

type userController struct {
}

var User userController

func (*userController) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.User.Register(&user); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("注册失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("注册成功", data))
	}
}

func (*userController) Login(c *gin.Context) {
	type loginForm struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var form loginForm

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("登录失败", err.Error()))
		return
	}

	if data, err := service.User.Login(form.Account, form.Password); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("登录失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("登录成功", data))
	}
}

func (*userController) Get(c *gin.Context) {
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

	if data, err := service.User.Get(page, name); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("获取用户信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("获取用户信息成功", data))
	}
}

func (*userController) Update(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	if data, err := service.User.Update(&user); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("修改用户信息失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("修改用户信息成功", data))
	}
}

func (*userController) Delete(c *gin.Context) {
	id := c.Query("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if data, err := service.User.Delete(oId); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("删除失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("删除成功", data))
	}
}

func (*userController) ResetPass(c *gin.Context) {
	id := c.Query("id")
	phone := c.Query("phone")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}
	if data, err := service.User.ResetPass(oId, phone); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("重置密码失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("重置密码成功", data))
	}
}

func (*userController) ChangePass(c *gin.Context) {
	type changePassForm struct {
		OldPass string `json:"oldPass" binding:"required"`
		NewPass string `json:"newPass" binding:"required"`
	}

	var form changePassForm

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("参数错误", err.Error()))
		return
	}

	id, ok := c.Get("Id")
	if ok == false {
		c.JSON(http.StatusOK, util.ErrorMess("Id错误", nil))
		return
	}
	if data, err := service.User.ChangePass(id, form.NewPass, form.OldPass); err != nil {
		c.JSON(http.StatusOK, util.ErrorMess("修改密码失败", err.Error()))
	} else {
		c.JSON(http.StatusOK, util.SuccessMess("修改密码成功", data))
	}
}
