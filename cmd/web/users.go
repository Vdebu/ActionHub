package main

import (
	"ActionHub/internal/data"
	"github.com/gin-gonic/gin"
)

func (app *application) registerUserHandler(c *gin.Context) {
	// 创建结构体用于读取用户输入的信息
	var input data.User
	// 使用gin的bind方法解析输入数据
	if err := c.ShouldBindJSON(&input); err != nil {
		// 报错并结束当前请求
		app.badRequestResponse(c, err)
		return
	}

}
