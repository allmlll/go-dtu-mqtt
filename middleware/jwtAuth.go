package middleware

import (
	"net/http"
	"ruitong-new-service/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = c.Request.URL.Query().Get("Authorization")
		}
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, util.ErrorMess("没有token", nil))
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == util.TokenPrefix) {
			c.JSON(http.StatusUnauthorized, util.ErrorMess("token格式错误", nil))
			c.Abort()
			return
		}
		tokenString := parts[1]
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		uc, err := util.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, util.ErrorMess("无效的Token", nil))
			c.Abort()
			return
		}
		// 将当前请求的Id信息保存到请求的上下文c上
		// 后续的处理函数可以用过c.PaginationFind("Id")来获取当前请求的用户Id
		c.Set("Id", uc.Id)
		c.Next()
	}
}
