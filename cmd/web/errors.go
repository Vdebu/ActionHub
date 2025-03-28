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

// 返回填写信息无效
func (app *application) invalidCredentialsResponse(c *gin.Context) {
	msg := "invalid authentication credentials"
	app.errorResponse(c, http.StatusUnauthorized, msg)
}

// 返回需要认证
func (app *application) authenticationRequiredResponse(c *gin.Context) {
	msg := "you must be authenticated to access this resource"
	app.errorResponse(c, http.StatusUnauthorized, msg)
}

// 返回表头存储的秘钥无效
func (app *application) invalidAuthenticationTokenResponse(c *gin.Context) {
	// 告诉客户端应该使用未加密的Token进行认证
	c.Header("WWW-Authenticate", "Bearer")
	msg := "invalid or missing authentication key"
	app.errorResponse(c, http.StatusUnauthorized, msg)
}
