package main

import (
	"ActionHub/internal/data"
	"github.com/gin-gonic/gin"
	"net/http"
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
	hashedPwd, err := app.hashPassword(input.Password)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 更新密码
	input.Password = hashedPwd
	// 生成JWT令牌
	token, err := app.generateJWT(input.Username)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 流程中未出现问题
	app.writeJSON(c, http.StatusOK, envelop{
		"JWT Token": token,
	})
}
