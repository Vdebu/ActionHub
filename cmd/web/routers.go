package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 初始化路由
func routers() *gin.Engine {
	// 使用gin的default(包含logger与recover)
	router := gin.Default()
	router.GET("/ping", func(context *gin.Context) {
		// 输出相应成功的信息
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return router
}
