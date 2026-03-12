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

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GeminiHandler — Upload resume, generate questions, link everything with an interview
func GeminiHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized — please login first"})
		return
	}

	jobTitle := c.PostForm("job_title")
	jobDescription := c.PostForm("job_description")
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please upload a resume file (PDF)"})
		return
	}

	// Convert user ID safely
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
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session data"})
		return
	}

	// Step 1: Create Interview record
	now := time.Now()
	interview := models.Interview{
		UserID:    uid,
		Status:    models.InterviewStatusActive, // ✅ fixed type-safe constant
		StartedAt: &now,
		JobTitle:  jobTitle, // 🆕 store job title directly
	}
	if err := config.DB.Create(&interview).Error; err != nil {
		fmt.Println("❌ Failed to create interview:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interview"})
		return
	}

	// Step 2: Ensure upload folder exists
	resumeDir := filepath.Join("uploads", "resumes")
	if err := os.MkdirAll(resumeDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resumes directory"})
		return
	}

	// Step 3: Save uploaded resume
	uniqueName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	savePath := filepath.Join(resumeDir, uniqueName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume file"})
		return
	}

	// Step 4: Extract text
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open resume"})
		return
	}
	defer f.Close()

	resumeText, err := services.ExtractResumeText(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to extract text: %v", err)})
		return
	}

	// Step 5: Save Resume record
	resume := models.Resume{
		UserID:       uid,
		InterviewID:  interview.ID,
		FilePath:     savePath,
		ExtractedTxt: resumeText,
	}
	if err := config.DB.Create(&resume).Error; err != nil {
		fmt.Println("❌ DB error while saving resume:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume"})
		return
	}

	// Step 6: Generate Questions using Gemini
	fmt.Printf("📨 Calling Gemini for Job Title: %s | Job Description length: %d\n", jobTitle, len(jobDescription))
	geminiResp, err := services.GenerateQuestions(jobTitle, jobDescription, resumeText, 5)
	if err != nil {
		fmt.Println("❌ Gemini error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(geminiResp.Questions) == 0 {
		fmt.Println("⚠️ No questions returned by Gemini API")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI did not return any questions"})
		return
	}

	// Step 7: Save Questions
	var savedQuestions []models.Question
	for _, q := range geminiResp.Questions {
		question := models.Question{
			InterviewID: interview.ID,
			ResumeID:    resume.ID,
			Text:        q.Text,
			IsAudio:     false,
			CreatedAt:   time.Now(),
		}
		if err := config.DB.Create(&question).Error; err != nil {
			fmt.Println("❌ Failed to insert question:", err)
			continue
		}
		savedQuestions = append(savedQuestions, question)
	}

	fmt.Printf("✅ %d questions saved for InterviewID=%d\n", len(savedQuestions), interview.ID)

	// Step 8: Return response
	c.JSON(http.StatusOK, gin.H{
		"message":          "Interview created, resume and questions saved ✅",
		"interview_id":     interview.ID,
		"resume_id":        resume.ID,
		"questions":        savedQuestions,
		"started_at":       interview.StartedAt,
		"interview_status": interview.Status,
	})
}
