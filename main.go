package main

import (
	"log"
	"os"

	"ai-backend-go/config"
	"ai-backend-go/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment & connect DB
	config.LoadEnv()
	config.ConnectDatabase()

	if os.Getenv("PORT") == "" {
		_ = os.Setenv("PORT", "8081")
	}
	port := os.Getenv("PORT")

	r := gin.Default()

	// CORS setup for frontend (React)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Cookie-based session
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("ai-session", store))

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "AI Backend Running ✅"})
	})

	// Route registration
	routes.RegisterAuthRoutes(r)
	routes.RegisterInterviewRoutes(r)
	routes.RegisterAIRoutes(r)
	routes.RegisterDashboardRoutes(r)

	// Launch server
	log.Printf("🚀 AI Backend running on http://127.0.0.1:%s", port)
	r.Run(":" + port)
}
