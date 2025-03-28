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

// 自动创建或更新表结构
func (m UserModel) Migrate(user *User) error {
	err := m.db.AutoMigrate(user)
	return err
}

// 向表中插入记录
func (m UserModel) Create(user *User) error {
	return m.db.Create(user).Error
}
