package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type envelop map[string]interface{}

// 输出错误信息
func (app *application) errorResponse(c *gin.Context, status int, message interface{}) {
	c.AbortWithStatusJSON(status, envelop{"error": message})
}

// 内部服务器错误
func (app *application) serverErrorResponse(c *gin.Context, err error) {
	// 在控制台输出错误的相关信息
	msg := "this server encountered a problem and can not process your request"
	app.errorResponse(c, http.StatusInternalServerError, msg)
}

// 用户端错误
func (app *application) badRequestResponse(c *gin.Context, err error) {
	// 输出错误相关信息
	app.errorResponse(c, http.StatusBadRequest, err.Error())
}
func (app *application) notFound(c *gin.Context) {
	msg := "the requested source can not be found"
	app.errorResponse(c, http.StatusNotFound, msg)
}
func (app *application) invalidCredentialsResponse(c *gin.Context) {
	msg := "invalid authentication credentials"
	app.errorResponse(c, http.StatusUnauthorized, msg)
}
