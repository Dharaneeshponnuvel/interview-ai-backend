package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"ai-backend-go/config"
	"ai-backend-go/models"

	"github.com/gin-gonic/gin"
)

func UploadResume(c *gin.Context) {
	db := config.DB
	userIDRaw, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	// interview_id is required
	interviewID := c.PostForm("interview_id")
	if interviewID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "interview_id is required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// save locally (or your S3/GCS storage)
	dst := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	res := models.Resume{
		UserID:      userID,
		InterviewID: parseUint(interviewID),
		FilePath:    dst,
	}
	if err := db.Create(&res).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resume_id":    res.ID,
		"interview_id": res.InterviewID,
		"file_path":    res.FilePath,
	})
}

func parseUint(s string) uint {
	var id uint
	_, _ = fmt.Sscan(s, &id)
	return id
}
