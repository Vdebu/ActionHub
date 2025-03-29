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
	// 权限不敏感的操作
	api.GET("/exchangeRate", app.getExchangeRate)
	api.GET("/articles", app.getArticles)
	api.GET("/articles/:id", app.getArticle)
	// 以下操作需要登入后才能进行
	api.Use(app.requireAuthentication())
	{
		// 使用POST方法创建新的汇率信息
		api.POST("/exchangeRate", app.createExchangeRate)
		// 创建新文章
		api.POST("/articles", app.createArticle)
		// 为文章点赞
		api.POST("/articles/:id/like", app.likeArticle)
	}
	return r
}
