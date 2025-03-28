package data

import "gorm.io/gorm"

type Models struct {
	Users UserModel
}

// 初始化数据模型
func NewModels(db *gorm.DB) Models {
	// 初始化数据库链接
	return Models{
		UserModel{db: db},
	}
}
