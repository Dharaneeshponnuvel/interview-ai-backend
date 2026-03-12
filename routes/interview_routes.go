package routes

import (
	"ai-backend-go/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterInterviewRoutes — all /interview endpoints
func RegisterInterviewRoutes(r *gin.Engine) {
	interview := r.Group("/interview")
	{
		// 🎯 Start a new interview (session + record)
		interview.POST("/start", handlers.StartInterview)

		// 🧠 Evaluate each answer (AI feedback)
		interview.POST("/evaluate", handlers.EvaluateHandler)

		// 🛑 End interview and save summary/score
		interview.POST("/:id/end", handlers.EndInterview)

		// 📊 Fetch complete results for dashboard
		interview.GET("/:id/results", handlers.GetInterviewResults)
		r.GET("/interviews/my", handlers.GetMyInterviews)
		r.GET("/interview/:id/report", handlers.DownloadInterviewReport)

	}
}
