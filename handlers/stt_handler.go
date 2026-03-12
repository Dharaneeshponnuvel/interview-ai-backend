package handlers

import (
	"fmt"
	"io"
	"net/http"

	"ai-backend-go/config"
	"ai-backend-go/models"
	"ai-backend-go/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// STTHandler — Converts speech to text and saves as answer (linked to interview/question)
func STTHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	interviewID := c.PostForm("interview_id")
	questionID := c.PostForm("question_id")
	if interviewID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "interview_id is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer f.Close()

	audioBytes, _ := io.ReadAll(f)
	transcript, err := services.SpeechToText(audioBytes, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized — please login first"})
		return
	}

	var uid, iID, qID uint
	fmt.Sscanf(interviewID, "%d", &iID)
	fmt.Sscanf(questionID, "%d", &qID)
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

	answer := models.Answer{
		UserID:      uid,
		InterviewID: iID,
		QuestionID:  &qID,
		AnswerText:  transcript,
		IsAudio:     true,
	}
	if err := config.DB.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transcript"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Transcription saved successfully",
		"transcript":   transcript,
		"answer_id":    answer.ID,
		"interview_id": iID,
		"question_id":  qID,
	})
}
