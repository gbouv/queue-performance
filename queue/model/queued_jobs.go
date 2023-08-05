package model

import (
	"time"
)

type QueuedJob struct {
	JobId       JobId     `gorm:"primaryKey"`
	CreatedTime time.Time `gorm:"index"`
	StartedTime *time.Time
	Difficulty  int
}
