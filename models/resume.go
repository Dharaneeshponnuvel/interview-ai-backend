package models

import "time"

type Resume struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `json:"user_id"`
	FilePath     string    `json:"file_path"`
	ExtractedTxt string    `gorm:"type:text" json:"extracted_text"`
	InterviewID  uint      `json:"interview_id" gorm:"index;not null"`
	JobType      string    `json:"job_type"`
	CreatedAt    time.Time `json:"created_at"`

	User      User       `gorm:"foreignKey:UserID"`
	Questions []Question `gorm:"constraint:OnDelete:CASCADE;"`
}
