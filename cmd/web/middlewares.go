package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) requireAuthentication() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 从请求头解析JWT token
		token := context.GetHeader("Authorization")
		if token == "" {
			// 返回表头存储的秘钥无效
			app.invalidAuthenticationTokenResponse(context)
			return
		}
		// 解析JWT
		username, err := app.parseJWT(token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(context)
			return
		}
		// 把读取到的用户名存入context
		context.Set("username", username)
		// 调用下一个中间件
		context.Next()
	}
}
