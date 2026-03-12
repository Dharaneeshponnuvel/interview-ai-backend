package routes

import (
	"ai-backend-go/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes — handles user authentication
func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/me", handlers.AuthMe)

	}
}
