package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ai-backend-go/config"
	"ai-backend-go/models"

	"github.com/gin-gonic/gin"
)

// GetInterviewResults — fetch final results for an interview (frontend-ready)
func GetInterviewResults(c *gin.Context) {
	idParam := c.Param("id")
	interviewID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interview ID"})
		return
	}

	var interview models.Interview
	if err := config.DB.First(&interview, interviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	// 🧑 Fetch user info
	var user models.User
	if err := config.DB.First(&user, interview.UserID).Error; err != nil {
		user = models.User{Name: "N/A", Email: "N/A"}
	}

	// 🧠 Fetch answers + related feedback + questions
	var answers []models.Answer
	if err := config.DB.
		Preload("Feedback").
		Preload("Question").
		Where("interview_id = ?", interviewID).
		Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch answers"})
		return
	}

	// 📊 Compute average score & build question list
	var totalScore float64
	var count int
	var questions []gin.H

	for _, ans := range answers {
		var score float64
		var feedback string
		if len(ans.Feedback) > 0 {
			score = ans.Feedback[0].Score
			feedback = ans.Feedback[0].Comments
			totalScore += score
			count++
		}
		questions = append(questions, gin.H{
			"id":       ans.QuestionID,
			"text":     ans.Question.Text,
			"score":    score,
			"feedback": feedback,
		})
	}

	avgScore := 0.0
	if count > 0 {
		avgScore = totalScore / float64(count)
	}

	// ⏱️ Duration
	var duration string
	if interview.StartedAt != nil && interview.EndedAt != nil {
		dur := interview.EndedAt.Sub(*interview.StartedAt)
		duration = dur.Round(time.Second).String()
	} else {
		duration = "N/A"
	}

	// ✅ Return structure that matches frontend fields exactly
	result := gin.H{
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
		},
		"job_title":     interview.JobTitle,
		"interview_id":  interview.ID,
		"status":        interview.Status,
		"started_at":    interview.StartedAt,
		"ended_at":      interview.EndedAt,
		"duration":      duration,
		"average_score": avgScore,
		"questions":     questions,
		"summary":       interview.Summary,
	}

	c.JSON(http.StatusOK, result)
}
