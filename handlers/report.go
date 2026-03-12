package handlers

import (
	"fmt"
	"net/http"

	"ai-backend-go/config"
	"ai-backend-go/models"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// ------------------------------
// DOWNLOAD INTERVIEW REPORT (PDF)
// ------------------------------
func DownloadInterviewReport(c *gin.Context) {
	interviewID := c.Param("id")

	// Fetch the interview
	var interview models.Interview
	if err := config.DB.First(&interview, "id = ?", interviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	// Fetch user
	var user models.User
	config.DB.First(&user, interview.UserID)

	// Fetch answers + feedback
	var answers []models.Answer
	config.DB.Preload("Feedback").
		Preload("Question").
		Where("interview_id = ?", interviewID).
		Find(&answers)

	// Calculate final score
	var finalScore float64
	config.DB.Table("feedbacks").
		Select("AVG(score)").
		Where("interview_id = ?", interviewID).
		Scan(&finalScore)

	// -------------------------------------------
	// CREATE PDF
	// -------------------------------------------
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(190, 10, "AI Interview Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	// ---------- Basic Information ----------
	pdf.MultiCell(0, 8, fmt.Sprintf("Interview ID: %d", interview.ID), "", "L", false)
	pdf.MultiCell(0, 8, fmt.Sprintf("Candidate: %s", user.Name), "", "L", false)
	pdf.MultiCell(0, 8, fmt.Sprintf("Email: %s", user.Email), "", "L", false)
	pdf.MultiCell(0, 8, fmt.Sprintf("Job Title: %s", interview.JobTitle), "", "L", false)
	pdf.MultiCell(0, 8, fmt.Sprintf("Final Score: %.2f", finalScore), "", "L", false)
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 14)
	pdf.MultiCell(0, 10, "Questions & Feedback", "", "L", false)
	pdf.SetFont("Arial", "", 12)

	// ---------- Questions, Answers, Feedback ----------
	for i, a := range answers {
		pdf.Ln(3)

		qHeader := fmt.Sprintf("Q%d: %s", i+1, a.Question.Text)
		pdf.MultiCell(0, 8, qHeader, "", "L", false)

		pdf.SetFont("Arial", "I", 12)
		pdf.MultiCell(0, 8, fmt.Sprintf("Answer: %s", a.AnswerText), "", "L", false)

		if len(a.Feedback) > 0 {
			pdf.SetFont("Arial", "", 12)
			pdf.MultiCell(0, 8, fmt.Sprintf("Score: %.2f", a.Feedback[0].Score), "", "L", false)
			pdf.MultiCell(0, 8, fmt.Sprintf("Feedback: %s", a.Feedback[0].Comments), "", "L", false)
		}

		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(0, 5, "-------------------------------------------", "", "L", false)
	}

	// -------------------------------------------
	// OUTPUT FILE
	// -------------------------------------------
	fileName := fmt.Sprintf("interview_report_%s.pdf", interviewID)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	err := pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
	}
}
