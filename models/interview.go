package models

import (
	"time"

	"gorm.io/gorm"
)

type InterviewStatus string

const (
	InterviewStatusCreated   InterviewStatus = "created"
	InterviewStatusActive    InterviewStatus = "active"
	InterviewStatusCompleted InterviewStatus = "completed"
)

type Interview struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	UserID       uint            `json:"user_id" gorm:"index;not null"`
	Status       InterviewStatus `json:"status" gorm:"type:varchar(20);default:'created';not null"`
	StartedAt    *time.Time      `json:"started_at"`
	EndedAt      *time.Time      `json:"ended_at"`
	Summary      *string         `json:"summary"`
	OverallScore *float64        `json:"overall_score"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index"`
	JobTitle     string          `json:"job_title" gorm:"type:varchar(255)"`
}
