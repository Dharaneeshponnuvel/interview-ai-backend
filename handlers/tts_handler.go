package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ai-backend-go/config"
	"ai-backend-go/models"
	"ai-backend-go/services"

	"github.com/gin-gonic/gin"
)

// TTSHandler — Convert text to speech and optionally link to interview/question
func TTSHandler(c *gin.Context) {
	var req struct {
		Text        string `json:"text" binding:"required"`
		Voice       string `json:"voice"`
		ReturnType  string `json:"return_type"`
		InterviewID uint   `json:"interview_id"`
		QuestionID  uint   `json:"question_id"`
		SaveToDB    bool   `json:"save_to_db"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Voice == "" {
		req.Voice = "en-US-Standard-C"
	}
	if req.ReturnType == "" {
		req.ReturnType = "file"
	}

	base64, audioBytes, err := services.TextToSpeech(req.Text, req.Voice, req.ReturnType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.SaveToDB && req.ReturnType == "file" {
		filename := fmt.Sprintf("tts_%d.mp3", time.Now().Unix())
		savePath := filepath.Join("uploads/audio", filename)
		os.MkdirAll(filepath.Dir(savePath), 0755)
		if err := os.WriteFile(savePath, audioBytes, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save audio"})
			return
		}

		audioRecord := models.Question{
			InterviewID: req.InterviewID,
			Text:        req.Text,
			IsAudio:     true,
			AudioPath:   savePath,
		}
		config.DB.Create(&audioRecord)
	}

	if req.ReturnType == "base64" {
		c.JSON(http.StatusOK, gin.H{"audio_base64": base64})
	} else {
		c.Data(http.StatusOK, "audio/mpeg", audioBytes)
	}
}
