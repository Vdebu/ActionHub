package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) writeJSON(c *gin.Context, code int, message interface{}) {
	c.IndentedJSON(code, message)
}
