package models

import (
	"time"
)

type Cache struct {
	Symbol    string `gorm:"primarykey"`
	StepSize  float64
	CreatedAt time.Time
}
