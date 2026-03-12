package models

import "time"

type Question struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InterviewID uint      `json:"interview_id" gorm:"index;not null"`
	ResumeID    uint      `json:"resume_id"`
	Text        string    `gorm:"type:text" json:"text"`
	IsAudio     bool      `json:"is_audio"`
	AudioPath   string    `json:"audio_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	Resume  Resume   `gorm:"foreignKey:ResumeID"`
	Answers []Answer `gorm:"constraint:OnDelete:CASCADE;"`
}
