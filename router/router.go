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

	r.GET("/", middleware.GetHome)
	r.GET("/products/:id", middleware.GetProductDetail)
	r.POST("/sign_in", middleware.LoginHandler)
	r.POST("/sign_up", middleware.RegisterHandler)
	r.POST("/profile", middleware.ProfileHandler)
	r.POST("/chat_lobby", middleware.GetChatLobby)
	r.GET("/search", middleware.GetSearch)
	r.GET("/ws", middleware.WebSocketHandler)

	return r
}
