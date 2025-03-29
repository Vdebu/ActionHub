package main

import (
	"ActionHub/internal/data"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"net/http"
	"time"
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
	// 创建新文章后将缓存中的旧数据删除(触发缓存未命中保证后续得到的是最新数据)
	// 上锁
	if err := app.mu.Lock(); err != nil {
		app.EditConflictResponse(c)
		return
	}
	// 删除操作
	if err := app.models.Articles.DelOldCache(); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 解锁
	if ok, err := app.mu.Unlock(); !ok || err != nil {
		panic("unlock failed")
		return
	}
	// 输出成功插入的数据
	app.writeJSON(c, http.StatusCreated, envelop{"article": article})
}

func (app *application) getArticles(c *gin.Context) {
	// 旁路缓存模式
	// 这里返回的是一个string(redis的硬性存储要求)后续要转为为[]byte反序列化
	cachedData, err := app.models.Articles.GetLatestInCache()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			// 如果返回的err是redis.nil则说明缓存不存在
			// 处理缓存未命中的情况
			var articles []data.Article
			if err := app.models.Articles.GetLatestInDatabase(&articles); err != nil {
				switch {
				case errors.Is(err, gorm.ErrRecordNotFound):
					app.notFound(c)
				default:
					app.serverErrorResponse(c, err)
				}
				return
			}
			// 将从数据库获取到的数据序列化为JSON存放进缓存中
			articleJSON, err := json.Marshal(articles)
			if err != nil {
				app.serverErrorResponse(c, err)
				return
			}
			// 使用对应键设置值(cachedArticle)
			if err := app.models.Articles.SetValueInCache(articleJSON, 10*time.Minute); err != nil {
				app.serverErrorResponse(c, err)
				return
			}
			// 输出查询到的内容
			app.writeJSON(c, http.StatusOK, envelop{"articles": articles})
			return
		default:
			app.serverErrorResponse(c, err)
			return
		}
	}
	// 缓存命中
	// 反序列化为原始数据
	var articles []data.Article
	if err := json.Unmarshal([]byte(cachedData), &articles); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 输出获取到的最新的文章
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
	// 使用分布式锁
	if err := app.mu.Lock(); err != nil {
		app.EditConflictResponse(c)
		return
	}
	// 将当前的键值+1
	if err := app.models.Articles.Likes(likeKey); err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	// 解锁
	if ok, err := app.mu.Unlock(); !ok || err != nil {
		panic("unlock failed")
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
