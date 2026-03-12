package handlers

import (
	"fmt"
	"net/http"
	"time"

	"ai-backend-go/config"
	"ai-backend-go/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func StartInterview(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var uid uint
	switch v := userID.(type) {
	case int:
		uid = uint(v)
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	case string:
		fmt.Sscanf(v, "%d", &uid)
	}

	// 🔥 FIX: Check if user already has a non-completed interview
	var existing models.Interview
	if err := config.DB.
		Where("user_id = ? AND status IN ?", uid, []string{
			string(models.InterviewStatusCreated),
			string(models.InterviewStatusActive),
		}).
		First(&existing).Error; err == nil {

		// Already has active/created interview
		c.JSON(http.StatusOK, gin.H{
			"message":      "Interview already running",
			"interview_id": existing.ID,
			"status":       existing.Status,
			"started_at":   existing.StartedAt,
		})
		return
	}

	// No ongoing interview → create new
	now := time.Now()
	newInterview := models.Interview{
		UserID:    uid,
		Status:    models.InterviewStatusActive,
		StartedAt: &now,
	}

	if err := config.DB.Create(&newInterview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interview"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Interview started",
		"interview_id": newInterview.ID,
		"status":       newInterview.Status,
		"started_at":   newInterview.StartedAt,
	})
}
func EndInterview(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	interviewID := c.Param("id")

	var interview models.Interview
	if err := config.DB.First(&interview, interviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	now := time.Now()
	interview.Status = models.InterviewStatusCompleted
	interview.EndedAt = &now

	// Calculate final score
	var finalAvg float64
	config.DB.
		Table("feedbacks").
		Select("AVG(score)").
		Where("interview_id = ?", interview.ID).
		Scan(&finalAvg)

	interview.OverallScore = &finalAvg

	if err := config.DB.Save(&interview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update interview"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Interview completed",
		"interview_id": interview.ID,
		"final_score":  finalAvg,
		"ended_at":     interview.EndedAt,
	})
}
