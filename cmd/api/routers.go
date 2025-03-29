package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// 初始化路由
func (app *application) routers() *gin.Engine {
	// 使用gin的default(包含logger与recover)
	r := gin.Default()
	// 配置CORS跨源请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9393"},                   // 允许的源
		AllowMethods:     []string{"GET", "POST"},                             // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 允许前端访问的响应头
		AllowCredentials: true,                                                // 允许发送验证凭证(cookie等)
		MaxAge:           12 * time.Hour,                                      // 预检请求的缓存时间
		//AllowOriginFunc: func(origin string) bool {}, 						  动态决定允许的源上边定义的切片就会失效
	}))
	auth := r.Group("/api/auth")
	{
		// 用户登入与注册
		auth.POST("/login", app.loginUserHandler)
		auth.POST("/register", app.registerUserHandler)
	}
	api := r.Group("/api")
	// 权限不敏感的操作
	api.GET("/exchangeRate", app.getExchangeRate)
	// 获取所有文章或通过id获取指定文章
	api.GET("/articles", app.getArticles)
	api.GET("/articles/:id", app.getArticle)
	// 获取指定文章的点赞信息
	api.GET("/articles/:id/like", app.getArticleLikes)
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
