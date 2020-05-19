package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/channels"
	"gitlab.com/kitalabs/go-2gaijin/middleware"
)

// Router is exported and used in main.go
func Router() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", middleware.GetHome)
	r.GET("/products/:id", middleware.GetProductDetail)
	r.GET("/wishlist", middleware.GetWishlistPage)

	r.POST("/sign_in", middleware.LoginHandler)
	r.POST("/sign_up", middleware.RegisterHandler)
	r.POST("/reset_password", middleware.ResetPasswordHandler)
	r.POST("/update_password", middleware.UpdatePasswordHandler)
	r.POST("/profile", middleware.ProfileHandler)
	r.POST("/confirm_identity", middleware.GenerateConfirmToken)
	r.GET("/confirm_email", middleware.EmailConfirmation)
	r.GET("/confirm_phone", middleware.PhoneConfirmation)

	r.POST("/chat_lobby", middleware.GetChatLobby)
	r.GET("/chat_messages", middleware.GetChatRoomMsg)
	r.GET("/initiate_chat", middleware.ChatUser)
	r.POST("/insert_message", middleware.InsertMessage)
	r.GET("/ws", channels.ServeChat)

	r.GET("/search", middleware.GetSearch)

	r.POST("/insert_notification", middleware.InsertNotification)
	r.POST("/insert_appointment", middleware.InsertAppointment)
	r.POST("/insert_trust_coin", middleware.InsertTrustCoin)

	r.GET("/get_appointments", middleware.GetAppointmentPage)
	r.GET("/get_notifications", middleware.GetNotificationPage)

	return r
}
