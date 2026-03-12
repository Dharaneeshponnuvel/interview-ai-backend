package models

import "time"

type Answer struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InterviewID uint      `json:"interview_id" gorm:"index;not null"`
	QuestionID  *uint     `json:"question_id"`
	UserID      uint      `json:"user_id"`
	AnswerText  string    `gorm:"type:text" json:"answer_text"`
	IsAudio     bool      `json:"is_audio"`
	AudioPath   string    `json:"audio_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	Question Question   `gorm:"foreignKey:QuestionID"`
	User     User       `gorm:"foreignKey:UserID"`
	Feedback []Feedback `gorm:"constraint:OnDelete:CASCADE;"`
}
