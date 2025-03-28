package data

import "gorm.io/gorm"

type User struct {
	gorm.Model        // 嵌入标准模型
	Username   string `gorm:"unique"` // 用户名唯一
	Password   string
}

type UserModel struct {
	db *gorm.DB
}

func (m UserModel) Insert() error {
	return nil
}
