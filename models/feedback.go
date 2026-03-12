package models

import "time"

type Feedback struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InterviewID uint      `json:"interview_id" gorm:"index;not null"`
	AnswerID    uint      `json:"answer_id"`
	Score       float64   `json:"score"`
	Comments    string    `gorm:"type:text" json:"comments"`
	CreatedAt   time.Time `json:"created_at"`

	Answer Answer `gorm:"foreignKey:AnswerID"`
}
