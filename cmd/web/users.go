package main

import (
	"ActionHub/internal/data"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) registerUserHandler(c *gin.Context) {
	// 创建结构体用于读取用户输入的信息
	var user data.User
	// 使用gin的bind方法解析输入数据
	if err := c.ShouldBindJSON(&user); err != nil {
		// 报错并结束当前请求
		app.badRequestResponse(c, err)
		return
	}
	hashedPwd, err := app.hashPassword(user.Password)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 更新密码
	user.Password = hashedPwd
	// 生成JWT令牌
	token, err := app.generateJWT(user.Username)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 更新或创建表结构
	err = app.models.Users.Migrate(&user)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 插入数据
	err = app.models.Users.Create(&user)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 流程中未出现问题
	app.writeJSON(c, http.StatusOK, envelop{
		"JWT Token": token,
	})
}
