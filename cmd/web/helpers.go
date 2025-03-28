package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
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

// 生成JWT令牌
func (app *application) generateJWT(username string) (string, error) {
	// 指定迁移方法与声明信息
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 声明信息
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // 过期时间(根据unix时间进行锚定)
	})
	// 获取完整签名的JWT
	// 填入加密与解密用秘钥
	signedToken, err := token.SignedString([]byte("secret"))
	// 方便操作直接加上Bearer前缀
	return "Bearer " + signedToken, err
}
