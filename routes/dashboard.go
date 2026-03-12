package routes

import (
	"ai-backend-go/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(r *gin.Engine) {
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/stats", handlers.DashboardStats)
		dashboard.GET("/score-history", handlers.ScoreHistory)
		dashboard.GET("/recent", handlers.DashboardRecent)
	}
}
