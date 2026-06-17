package models

import "time"

type Job struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Script    string    `json:"script" gorm:"not null"`
	Schedule  string    `json:"schedule"`
	Enabled   bool      `json:"enabled" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Run struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	JobID      uint       `json:"job_id" gorm:"index;not null"`
	Job        *Job       `json:"job,omitempty" gorm:"foreignKey:JobID"`
	Status     string     `json:"status" gorm:"default:pending"`
	Logs       string     `json:"logs" gorm:"type:text"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
