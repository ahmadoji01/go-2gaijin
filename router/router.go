package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/middleware"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/api/", middleware.GetHome)
	r.GET("/api/products", middleware.GetAllProducts)
	r.POST("/api/login", middleware.LoginHandler)
	r.POST("/api/register", middleware.RegisterHandler)
	r.POST("/api/profile", middleware.ProfileHandler)
	r.GET("/api/search", middleware.GetSearch)
	r.GET("/api/graphql", middleware.GetTestGraphQL)
	r.GET("/ws", middleware.WebSocketHandler)

	return r
}
