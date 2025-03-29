package main

import (
	"ActionHub/internal/data"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
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

// 解析传入的包含JWT的请求头kv
func (app *application) parseJWT(tokenString string) (string, error) {
	// 解析头部存储的秘钥信息
	tokenParts := strings.Split(tokenString, " ")
	if tokenParts[0] != "Bearer" || len(tokenParts) < 2 {
		return "", nil
	}
	token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
		// 检测方法是否相同
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// 返回加密/解密用秘钥
		return []byte("secret"), nil
	})
	// 判断是否发生错误
	if err != nil {
		return "", err
	}
	// 检验细节是否正确
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username must be a string")
		}
		return username, nil
	}
	return "", err
}

// 读取url中的参数
func (app *application) readID(c *gin.Context) (string, error) {
	stringID := c.Param("id")
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		return "", err
	}
	// 输入的查询参数无效
	if stringID == "" || ID < 1 {
		return "", data.ErrInvalidIDParams
	}
	return stringID, nil
}
