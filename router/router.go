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
	r.POST("/add_product", middleware.PostNewProduct)
	r.POST("/get_product_info_edit", middleware.GetProductInfoForEdit)
	r.POST("/mark_as_sold", middleware.MarkAsSold)
	r.POST("/edit_pricing", middleware.EditPricing)
	r.POST("/delete_product", middleware.DeleteProduct)
	r.POST("/edit_product", middleware.EditProduct)
	r.POST("/like_product", middleware.LikeProduct)
	r.GET("/get_categories", middleware.GetAllCategories)

	r.POST("/sign_in", middleware.LoginHandler)
	r.POST("/sign_up", middleware.RegisterHandler)
	r.POST("/sign_out", middleware.LogoutHandler)
	r.POST("/refresh_token", middleware.RefreshToken)
	r.POST("/reset_password", middleware.ResetPasswordHandler)
	r.POST("/update_password", middleware.UpdatePasswordHandler)
	r.POST("/profile", middleware.ProfileHandler)
	r.POST("/get_profile_info", middleware.GetProfileInfo)
	r.POST("/update_profile", middleware.UpdateProfile)
	r.POST("/upload_avatar", middleware.UploadProfilePhoto)
	r.POST("/confirm_identity", middleware.GenerateConfirmToken)
	r.GET("/check_notif_read", middleware.CheckNotifRead)
	r.GET("/confirm_email", middleware.EmailConfirmation)
	r.GET("/confirm_phone", middleware.PhoneConfirmation)

	r.GET("/profile_visitor", middleware.GetProfileForVisitorPage)

	r.POST("/chat_lobby", middleware.GetChatLobby)
	r.GET("/chat_messages", middleware.GetChatRoomMsg)
	r.GET("/chat_users", middleware.GetChatRoomUser)
	r.GET("/initiate_chat", middleware.ChatUser)
	r.POST("/insert_message", middleware.InsertMessage)
	r.POST("/insert_image_message", middleware.InsertImageMessage)

	r.GET("/ws", channels.ServeChat)

	r.GET("/search", middleware.GetSearch)

	r.POST("/insert_notification", middleware.InsertNotification)
	r.POST("/insert_appointment", middleware.InsertAppointment)
	r.POST("/insert_trust_coin", middleware.InsertTrustCoin)
	r.POST("/confirm_appointment", middleware.AppointmentConfirmation)
	r.POST("/reschedule_appointment", middleware.RescheduleAppointment)
	r.POST("/finish_appointment", middleware.FinishAppointment)
	r.GET("/get_seller_appointments", middleware.GetSellerAppointmentPage)
	r.GET("/get_buyer_appointments", middleware.GetBuyerAppointmentPage)
	r.GET("/get_notifications", middleware.GetNotificationPage)

	// Preflight Response
	r.OPTIONS("/", middleware.HandlePreflight)
	r.OPTIONS("/products/:id", middleware.HandlePreflight)
	r.OPTIONS("/wishlist", middleware.HandlePreflight)
	r.OPTIONS("/mark_as_sold", middleware.HandlePreflight)
	r.OPTIONS("/delete_product", middleware.HandlePreflight)
	r.OPTIONS("/edit_product", middleware.HandlePreflight)
	r.OPTIONS("/like_product", middleware.HandlePreflight)
	r.OPTIONS("/get_categories", middleware.HandlePreflight)
	r.OPTIONS("/sign_in", middleware.HandlePreflight)
	r.OPTIONS("/sign_up", middleware.HandlePreflight)
	r.OPTIONS("/sign_out", middleware.HandlePreflight)
	r.OPTIONS("/refresh_token", middleware.HandlePreflight)
	r.OPTIONS("/reset_password", middleware.HandlePreflight)
	r.OPTIONS("/update_password", middleware.HandlePreflight)
	r.OPTIONS("/profile", middleware.HandlePreflight)
	r.OPTIONS("/update_profile", middleware.HandlePreflight)
	r.OPTIONS("/confirm_identity", middleware.HandlePreflight)
	r.OPTIONS("/confirm_email", middleware.HandlePreflight)
	r.OPTIONS("/confirm_phone", middleware.HandlePreflight)
	r.OPTIONS("/profile_visitor", middleware.HandlePreflight)
	r.OPTIONS("/chat_lobby", middleware.HandlePreflight)
	r.OPTIONS("/chat_messages", middleware.HandlePreflight)
	r.OPTIONS("/initiate_chat", middleware.HandlePreflight)
	r.OPTIONS("/insert_message", middleware.HandlePreflight)
	r.OPTIONS("/search", middleware.HandlePreflight)
	r.OPTIONS("/insert_notification", middleware.HandlePreflight)
	r.OPTIONS("/insert_appointment", middleware.HandlePreflight)
	r.OPTIONS("/insert_trust_coin", middleware.HandlePreflight)
	r.OPTIONS("/confirm_appointment", middleware.HandlePreflight)
	r.OPTIONS("/reschedule_appointment", middleware.HandlePreflight)
	r.OPTIONS("/finish_appointment", middleware.HandlePreflight)
	r.OPTIONS("/get_seller_appointments", middleware.HandlePreflight)
	r.OPTIONS("/get_buyer_appointments", middleware.HandlePreflight)
	r.OPTIONS("/get_notifications", middleware.HandlePreflight)
	r.OPTIONS("/check_notif_read", middleware.HandlePreflight)
	r.OPTIONS("/chat_users", middleware.HandlePreflight)
	r.OPTIONS("/upload_avatar", middleware.HandlePreflight)
	r.OPTIONS("/add_product", middleware.HandlePreflight)
	r.OPTIONS("/get_product_info_edit", middleware.HandlePreflight)
	r.OPTIONS("/edit_pricing", middleware.HandlePreflight)
	r.OPTIONS("/get_profile_info", middleware.HandlePreflight)
	r.OPTIONS("/insert_image_message", middleware.HandlePreflight)

	return r
}
