package main

import (
	"ActionHub/internal/data"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"net/http"
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
	stringID, err := app.readID(c)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidIDParams):
			app.badRequestResponse(c, data.ErrInvalidIDParams)
		default:
			app.serverErrorResponse(c, err)
		}
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
func (app *application) likeArticle(c *gin.Context) {
	// 读取参数
	stringID, err := app.readID(c)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidIDParams):
			app.badRequestResponse(c, data.ErrInvalidIDParams)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}
	// 定义当前文章的article:ID:likes(redis中的kv规范)
	// 使用字符串拼接进行定义
	likeKey := "article:" + stringID + ":likes"
	// 将当前的键值+1
	if err := app.models.Articles.Likes(likeKey); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 返回成功点赞的信息
	app.writeJSON(c, http.StatusOK, envelop{
		"message": "Successfully liked the article",
	})
}

// 获取点赞数
func (app *application) getArticleLikes(c *gin.Context) {
	stringID, err := app.readID(c)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidIDParams):
			app.badRequestResponse(c, data.ErrInvalidIDParams)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}
	likeKey := "article:" + stringID + ":likes"
	// 获取redis中的键值
	likes, err := app.models.Articles.GetLikes(likeKey)
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			// 若当前点赞的信息不存在则表示当前点赞为0
			likes = "0"
			// 不需要return
		default:
			app.serverErrorResponse(c, err)
			return
		}
	}
	// 返回获取到的数据
	app.writeJSON(c, http.StatusOK, envelop{"likes": likes})
}
