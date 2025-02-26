package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"

	"net/http"
)

func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		method := c.Request.Method
		// 查找api
		var api model.Api
		if err := global.ApiColl.FindOne(context.TODO(), bson.M{"url": url, "method": method}).Decode(&api); err != nil {
			c.JSON(http.StatusForbidden, util.ErrorMess("验证api：此api不存在", err.Error()))
			c.Abort()
			return
		}

		id, _ := c.Get("Id")
		// 查找用户
		var user model.User
		err := global.UserColl.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusForbidden, util.ErrorMess("验证api：获取用户失败", err.Error()))
			c.Abort()
			return
		}

		// 查找角色
		var role model.Role
		err = global.RoleColl.FindOne(context.TODO(), bson.M{"_id": user.Role}).Decode(&role)
		if err != nil {
			c.JSON(http.StatusForbidden, util.ErrorMess("验证api：获取用户角色失败", err.Error()))
			c.Abort()
			return
		}
		//轮询role对应的apis，判断其是否相应的权限
		for _, api = range role.Apis {
			if api.Url == url && api.Method == method {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, util.ErrorMess("验证api：此用户无访问此api的权限", nil))
		c.Abort()
		return
	}
}
