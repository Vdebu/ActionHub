package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) writeJSON(c *gin.Context, code int, message interface{}) {
	c.IndentedJSON(code, message)
}

// 将传入的密码进行哈希
func (app *application) hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	// 没有后续操作可以直接返回
	return string(hash), err
}
