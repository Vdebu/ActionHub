package main

import (
	"ActionHub/internal/data"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (app *application) createExchangeRate(c *gin.Context) {
	var exchangeRate data.ExchangeRate
	// 尝试读取数据
	if err := c.ShouldBindJSON(&exchangeRate); err != nil {
		app.badRequestResponse(c, err)
		return
	}
	// 初始化时间字段
	exchangeRate.Date = time.Now().Local()
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
func (app *application) getExchangeRate(c *gin.Context) {
	// 存储查询到的汇率信息
	var exchangeRates []*data.ExchangeRate
	// 从数据库中读取数据
	if err := app.models.ExchangeRate.GetLatest(&exchangeRates); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 输出查找到的信息
	app.writeJSON(c, http.StatusOK, envelop{"exchange_rates": exchangeRates})
}
