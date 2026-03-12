package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"ai-backend-go/config"
	"ai-backend-go/models"
	"ai-backend-go/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func EvaluateHandler(c *gin.Context) {
	var req struct {
		InterviewID uint   `json:"interview_id" binding:"required"`
		JobTitle    string `json:"job_title"`
		QuestionID  uint   `json:"question_id"`
		Question    string `json:"question"`
		Answer      string `json:"answer"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

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

	// ---------------------------
	// AI Evaluation
	// ---------------------------
	fmt.Printf("🧠 Evaluating answer: %s\n", req.Answer)

	resp, err := services.EvaluateAnswer(req.JobTitle, req.Question, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI evaluation failed"})
		return
	}

	score := float64(resp.Score)

	// ---------------------------
	// Save Answer
	// ---------------------------
	var qid *uint
	if req.QuestionID != 0 {
		qid = &req.QuestionID
	}

	answer := models.Answer{
		UserID:      uid,
		InterviewID: req.InterviewID,
		QuestionID:  qid,
		AnswerText:  req.Answer,
	}

	config.DB.Create(&answer)

	// ---------------------------
	// Save Feedback
	// ---------------------------
	feedback := models.Feedback{
		InterviewID: req.InterviewID,
		AnswerID:    answer.ID,
		Score:       score,
		Comments:    strings.Join(resp.Improvements, ", "),
	}

	config.DB.Create(&feedback)

	// ---------------------------
	// 🆕 UPDATE INTERVIEW SCORE
	// ---------------------------
	var avgScore float64
	config.DB.
		Table("feedbacks").
		Select("AVG(score)").
		Where("interview_id = ?", req.InterviewID).
		Scan(&avgScore)

	config.DB.Model(&models.Interview{}).
		Where("id = ?", req.InterviewID).
		Update("overall_score", avgScore)

	// ---------------------------
	// RESPONSE
	// ---------------------------
	c.JSON(http.StatusOK, gin.H{
		"interview_id":    req.InterviewID,
		"question_id":     req.QuestionID,
		"answer_id":       answer.ID,
		"feedback_id":     feedback.ID,
		"score":           resp.Score,
		"strengths":       resp.Strengths,
		"improvements":    resp.Improvements,
		"sample_response": resp.SampleResponse,
	})
}
