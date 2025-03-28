package data

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type User struct {
	gorm.Model        // 嵌入标准模型
	Username   string `gorm:"unique"` // 用户名唯一
	Password   string
}

type UserModel struct {
	db *gorm.DB
}

// 声明User结构体对应操作的数据库中的表
func (u User) TableName() string {
	// 显式指定要操作的表名
	return "users"
}

// 自动创建或更新表结构
func (m UserModel) Migrate(user *User) error {
	err := m.db.AutoMigrate(user)
	return err
}

// 向表中插入记录
func (m UserModel) Create(user *User) error {
	// 这里也同样会返回自动生成的信息如id与创建时间(但只取错误信息)
	return m.db.Create(user).Error
}

// 通过用户名查找用户
func (m UserModel) GetByName(username string, user *User) error {
	// 查询排序后的第一条记录并将记录存储进dst
	return m.db.Where("username = ?", username).First(user).Error
}

// 比较用户输入的密码
func (m UserModel) Check(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	// 如果密码一直的话err就会是nil
	return err == nil
}
