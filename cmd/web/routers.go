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
		auth.POST("/login", app.loginUserHandler)
		auth.POST("/register", app.registerUserHandler)
	}
	api := r.Group("/api")
	// 获取最新的汇率信息是不需要登入权限的
	api.GET("/exchangeRate", app.getExchangeRate)
	// 以下操作需要登入后才能进行
	api.Use(app.requireAuthentication())
	{
		// 使用POST方法创建新的汇率信息
		api.POST("/exchangeRate", app.createExchangeRate)
	}
	return r
}
