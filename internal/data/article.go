package data

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrInvalidIDParams = errors.New("invalid id params")
)

type Article struct {
	gorm.Model
	Title   string `binding:"required"` // 使用binding required就不需要额外的json tag
	Content string `binding:"required"`
	Preview string `binding:"required"`
	Likes   int    `gorm:"default:0"` // 默认初始化为0
}

type ArticleModel struct {
	db *gorm.DB
}

// 显式指定表名
func (a Article) TableName() string {
	return "articles"
}

// 迁移数据库
func (m ArticleModel) Migrate(article *Article) error {
	return m.db.AutoMigrate(article)
}

// 创建新文章
func (m ArticleModel) Create(article *Article) error {
	return m.db.Create(article).Error
}

// 获取最新文章用于展示
func (m ArticleModel) GetLatest(articles *[]Article) error {
	return m.db.Find(articles).Error
}

// 根据提供的id查询文章
func (m ArticleModel) GetByID(id string, article *Article) error {
	return m.db.Where("id = ?", id).First(article).Error
}
