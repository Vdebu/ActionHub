package data

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Models struct {
	Users        UserModel
	ExchangeRate ExchangeRateModel
	Articles     ArticleModel
}

// 初始化数据模型
func NewModels(db *gorm.DB, redisDB *redis.Client) Models {
	// 初始化数据库链接
	return Models{
		Users:        UserModel{db: db},
		ExchangeRate: ExchangeRateModel{db: db},
		Articles:     ArticleModel{db: db, redisDB: redisDB, cacheKey: "cachedArticles"},
	}
}
