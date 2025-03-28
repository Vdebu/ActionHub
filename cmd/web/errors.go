package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 输出错误信息
func (app *application) errorResponse(c *gin.Context, status int, message interface{}) {
	c.AbortWithStatusJSON(status, message)
}

// 内部服务器错误
func (app *application) serverErrorResponse(c *gin.Context, err error) {
	// 在控制台输出错误的相关信息

	msg := "this server encountered a problem and can not process your request"
	app.errorResponse(c, http.StatusInternalServerError, msg)
}
