package routes

import (
	"ai-backend-go/handlers"

	"github.com/gin-gonic/gin"
)

// Register AIRoutes — handles general AI-related endpoints
func RegisterAIRoutes(r *gin.Engine) {
	ai := r.Group("/ai")
	{
		// ✳️ Generates interview questions using Gemini
		ai.POST("/generate-questions", handlers.GeminiHandler)

		// 🧠 Evaluates an individual question response
		ai.POST("/evaluate", handlers.EvaluateHandler)

	}
}
