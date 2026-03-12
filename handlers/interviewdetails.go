package handlers

import (
	"ai-backend-go/config"
	"ai-backend-go/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetMyInterviews(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var interviews []models.Interview

	if err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&interviews).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load interviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviews})
}
