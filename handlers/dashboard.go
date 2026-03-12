package handlers

import (
	"net/http"
	"time"

	"ai-backend-go/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ----------------------------
// 1. GET /dashboard/stats
// ----------------------------
func DashboardStats(c *gin.Context) {

	// ----------------------------
	// ✅ Get Logged-in User ID
	// ----------------------------
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
	}

	// Total interviews (for this user)
	var totalInterviews int64
	config.DB.Table("interviews").
		Where("user_id = ?", uid).
		Count(&totalInterviews)

	// Average score for this user
	var avg *float64
	config.DB.
		Table("interviews").
		Select("AVG(COALESCE(overall_score, 0))").
		Where("user_id = ?", uid).
		Scan(&avg)

	c.JSON(http.StatusOK, gin.H{
		"total_interviews": totalInterviews,
		"average_score":    avg,
		"user_id":          uid,
	})
}

// ----------------------------
// 2. GET /dashboard/score-history
// ----------------------------
func ScoreHistory(c *gin.Context) {

	// Get logged-in user ID
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
	}

	type Row struct {
		Date  time.Time
		Avg   float64
		Count int64
	}

	var rows []Row

	end := time.Now()
	start := end.AddDate(0, 0, -30)

	// Only logged-in user data
	config.DB.
		Table("interviews").
		Select("DATE(created_at) AS date, AVG(COALESCE(overall_score, 0)) AS avg, COUNT(*) AS count").
		Where("created_at BETWEEN ? AND ? AND user_id = ?", start, end, uid).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC").
		Scan(&rows)

	type Out struct {
		Date  string  `json:"date"`
		Avg   float64 `json:"avg"`
		Count int64   `json:"count"`
	}

	var out []Out
	for _, r := range rows {
		out = append(out, Out{
			Date:  r.Date.Format("2006-01-02"),
			Avg:   r.Avg,
			Count: r.Count,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": out})
}

// ----------------------------
// 3. GET /dashboard/recent
// ----------------------------
func DashboardRecent(c *gin.Context) {

	// Get user session
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
	}

	type Row struct {
		ID        uint     `json:"id"`
		UserName  string   `json:"user_name"`
		JobTitle  string   `json:"job_title"`
		CreatedAt string   `json:"created_at"`
		Score     *float64 `json:"score"`
	}

	var results []Row

	// Only logged-in user's interviews
	config.DB.
		Table("interviews").
		Select(`
			interviews.id,
			users.name AS user_name,
			interviews.job_title,
			interviews.created_at,
			interviews.overall_score
		`).
		Joins("LEFT JOIN users ON users.id = interviews.user_id").
		Where("interviews.user_id = ?", uid).
		Order("interviews.created_at DESC").
		Limit(10).
		Scan(&results)

	// Format Date
	for i, row := range results {
		t, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
		if err == nil {
			results[i].CreatedAt = t.Format("2006-01-02")
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
