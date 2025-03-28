package main

import (
	"github.com/gin-gonic/gin"
)

// 初始化路由
func (app *application) routers() *gin.Engine {
	// 使用gin的default(包含logger与recover)
	r := gin.Default()
	auth := r.Group("/api/auth")
	{
		// 用户登入与注册
		auth.POST("/login", app.registerUserHandler)
		auth.POST("/register")
	}
	return r
}
