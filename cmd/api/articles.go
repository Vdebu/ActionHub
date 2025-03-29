package main

import (
	"ActionHub/internal/data"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (app *application) createArticle(c *gin.Context) {
	var article data.Article
	// 读取数据
	if err := c.ShouldBindJSON(&article); err != nil {
		app.badRequestResponse(c, err)
		return
	}
	// 数据读取成功
	// 先进行迁移
	if err := app.models.Articles.Migrate(&article); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 插入数据
	if err := app.models.Articles.Create(&article); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 输出成功插入的数据
	app.writeJSON(c, http.StatusCreated, envelop{"article": article})
}

func (app *application) getArticles(c *gin.Context) {
	var articles []data.Article
	// 获取最新的文章展示
	if err := app.models.Articles.GetLatest(&articles); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 输出最新的文章
	app.writeJSON(c, http.StatusOK, envelop{"articles": articles})
}

// 通过id获取文章
func (app *application) getArticle(c *gin.Context) {
	stringID := c.Param("id")
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 输入的查询参数无效
	if stringID == "" || ID < 1 {
		app.badRequestResponse(c, data.ErrInvalidIDParams)
		return
	}
	// 根据解析出来的参数查询相关数据
	var article data.Article
	if err := app.models.Articles.GetByID(stringID, &article); err != nil {
		// 根据错误类型决定使用的状态码
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFound(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}
	// 输出生成的结果
	app.writeJSON(c, http.StatusOK, envelop{"article": article})
}
