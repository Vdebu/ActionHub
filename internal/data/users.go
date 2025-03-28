package data

import "gorm.io/gorm"

type User struct {
}

type UserModel struct {
	db *gorm.DB
}
