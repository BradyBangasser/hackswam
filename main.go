package main

import (
	"hackswam/m/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	routes.CreateRouter(engine)
	engine.Run()
}
