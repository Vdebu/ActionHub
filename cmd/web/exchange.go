package main

import (
	"ActionHub/internal/data"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) createExchangeRate(c *gin.Context) {
	var exchangeRate data.ExchangeRate
	// 尝试读取数据
	if err := c.ShouldBindJSON(&exchangeRate); err != nil {
		app.badRequestResponse(c, err)
		return
	}
	// 迁移数据库
	if err := app.models.ExchangeRate.Migrate(&exchangeRate); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 创建好表后进行记录的创建
	if err := app.models.ExchangeRate.Create(&exchangeRate); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 输出创建成功的信息
	app.writeJSON(c, http.StatusCreated, envelop{"exchange_rate": exchangeRate})
}
