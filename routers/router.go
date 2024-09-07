package routers

import (
	"IMProject/middleware"
	"github.com/gin-gonic/gin"
)

var Routes *gin.Engine

func InitRouter() {
	gin.ForceConsoleColor()
	Routes = gin.Default()
	middleware.InitJaeger()
	Routes.Use(middleware.Jaeger())
}
