package hello

import "github.com/gin-gonic/gin"

func GET(c *gin.Context) {
	c.String(200, "Hello World!")
}
