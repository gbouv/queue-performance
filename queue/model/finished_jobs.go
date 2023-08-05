package model

import (
	"time"
)

type FinishedJob struct {
	JobId       JobId `gorm:"primaryKey"`
	StartedTime time.Time
	DurationMs  uint
}
