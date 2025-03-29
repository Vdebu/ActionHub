package main

import (
	"ActionHub/internal/data"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			app.duplicateKeyResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}
	// 流程中未出现问题
	app.writeJSON(c, http.StatusCreated, envelop{
		"JWT Token": token,
	})
}
func (app *application) loginUserHandler(c *gin.Context) {
	// 序列化与反序列化tag
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// 读取用户输入的内容
	if err := c.ShouldBindJSON(&input); err != nil {
		app.badRequestResponse(c, err)
		return
	}
	// 解析成功载入user模型
	var user data.User
	// 使用gorm内置方法进行行级比对(先从用户名下手查看是否存在(unique))
	// 查询排序后的第一条结果并将查找到的结果存储进dst(如果有的话)
	err := app.models.Users.GetByName(input.Username, &user)
	if err != nil {
		// http 401
		app.invalidCredentialsResponse(c)
		return
	}
	// 检测密码是否正确
	if !app.models.Users.Check(input.Password, user.Password) {
		app.invalidCredentialsResponse(c)
		return
	}
	// 登录成功返回JWT给客户端
	token, err := app.generateJWT(user.Username)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 输出创建成功的信息并返回生成的JWT
	app.writeJSON(c, http.StatusOK, envelop{"token": token})
}
